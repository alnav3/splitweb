package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/alnav3/splitweb/auth"
	"github.com/alnav3/splitweb/templates"
)

// Profile page handler
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := getMockUserProfile()
	err := templates.ProfilePage(user).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

// Profile form handlers
func ProfileChangeEmailHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	newEmail := r.FormValue("new-email")
	currentPassword := r.FormValue("current-password")

	// In a real app, validate password and update email
	log.Printf("Email change requested: new email: %s, password provided: %t", newEmail, currentPassword != "")

	// Return success response - the frontend will handle the popup
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func ProfileChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Get user from session
	token, userId, authenticated := auth.GetUserFromSession(r)
	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please log in to change your password"))
		return
	}

	// Get user email for re-login
	userRecord, err := auth.GetUserRecord(token, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve user information"))
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
	authResponse, err := auth.AuthWithPassword(userRecord.Email, newPassword)
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
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a real app, delete user account and all associated data
	log.Printf("Account deletion requested")

	// For demo purposes, redirect to login
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

// Helper function to create mock user profile data
func getMockUserProfile() templates.UserProfile {
	return templates.UserProfile{
		ID:            "user123",
		Name:          "John Doe",
		Email:         "john.doe@example.com",
		CreatedAt:     time.Now().AddDate(0, -8, -15), // 8 months and 15 days ago
		LastLoginAt:   time.Now().Add(-time.Hour * 3), // 3 hours ago
		AvatarURL:     "https://ui-avatars.com/api/?name=John+Doe&background=89b4fa&color=fff",
		GroupCount:    5,
		TotalExpenses: 1247.50,
	}
}
