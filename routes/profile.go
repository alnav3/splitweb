package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alnav3/splitweb/auth"
	"github.com/alnav3/splitweb/templates"
)

// Profile page handler
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	_, userId, isValid := auth.GetUserFromSession(r)
	if isValid {
		user, _ := Repo.Queries.FindUserById(Repo.Context, userId) // add error handling later
		groupsCount, _ := Repo.Queries.GroupsCountByUserId(Repo.Context, userId) // add error handling later
		err := templates.ProfilePage(user, int(groupsCount)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
		return
	}
	http.Error(w, "Error rendering template", http.StatusInternalServerError)
	log.Printf("Error rendering template: userId not valid")
}

// verifyUserAndPassword validates session and current password
func verifyUserAndPassword(r *http.Request, currentPassword string) (token, userId, userEmail string, err error) {
	// Get user from session
	token, userId, authenticated := auth.GetUserFromSession(r)
	if !authenticated {
		return "", "", "", fmt.Errorf("please log in to perform this action")
	}

	// Get user email for password verification
	userRecord, err := auth.GetUserRecord(token, userId)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to retrieve user information")
	}

	// Verify current password
	err = auth.VerifyCurrentPassword(userRecord.Email, currentPassword)
	if err != nil {
		return "", "", "", fmt.Errorf("incorrect current password")
	}

	return token, userId, userRecord.Email, nil
}

// Profile form handlers
func ProfileChangeEmailHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	newEmail := r.FormValue("new-email")

	// Basic validation
	if newEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("New email is required"))
		return
	}

	token, _, authenticated := auth.GetUserFromSession(r)
	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please log in to perform this action"))
		return
	}

	// Request email change
	err = auth.RequestEmailChange(newEmail, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to request email change: " + err.Error()))
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func ProfileChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current-password")
	newPassword := r.FormValue("new-password")
	confirmPassword := r.FormValue("confirm-new-password")

	// Basic validation
	if newPassword != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords do not match"))
		return
	}

	if len(newPassword) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password must be at least 8 characters long"))
		return
	}

	// Verify user and password
	token, userId, userEmail, err := verifyUserAndPassword(r, currentPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Change password
	err = auth.ChangePasswordWithOlderOne(currentPassword, newPassword, confirmPassword, token, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to change password: " + err.Error()))
		return
	}

	// Re-authenticate with new password
	authResponse, err := auth.AuthWithPassword(userEmail, newPassword, Repo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Password changed but failed to re-authenticate. Please log in again."))
		return
	}

	// Update session with new token
	err = auth.SetUserSession(r, w, authResponse.Token, authResponse.Record.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Password changed but failed to update session. Please log in again."))
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func ProfileDeleteAccountHandler(w http.ResponseWriter, r *http.Request) {

	currentPassword := r.FormValue("current-password")

	// Basic validation
	if currentPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Current password is required to delete account"))
		return
	}

	// Verify user and password
	token, userId, _, err := verifyUserAndPassword(r, currentPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Delete user from PocketBase
	err = auth.DeleteUserAccount(token, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to delete account: " + err.Error()))
		return
	}

	// Delete user from local database
	err = Repo.Queries.DeleteUserById(Repo.Context, userId)
	if err != nil {
		log.Printf("Failed to delete user from local database: %v", err)
		// Continue even if local deletion fails, as the PocketBase deletion succeeded
	}

	// Clear the session and redirect to login
	err = auth.ClearSession(r, w)
	if err != nil {
		log.Printf("Failed to clear session: %v", err)
		// Continue anyway, as the account deletion succeeded
	}

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

