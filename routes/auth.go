package routes

import (
	"log"
	"net/http"

	"github.com/alnav3/splitweb/auth"
	"github.com/alnav3/splitweb/templates"
)

// Authentication page handlers
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.Login().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.Register().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ForgotPassword().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

// Authentication form handlers
func AuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Basic validation
	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("Email and password are required.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering validation error: %v", err)
		}
		return
	}

	// Authenticate with PocketBase
	authResponse, err := auth.AuthWithPassword(email, password, Repo)
	if err != nil {
		log.Printf("Authentication failed for %s: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		err := templates.LoginErrorPopup().Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering login error popup: %v", err)
		}
		return
	}

	// Set user session
	err = auth.SetUserSession(r, w, authResponse.Token, authResponse.Record.Id)
	if err != nil {
		log.Printf("Error setting user session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		err := templates.AuthValidationError("Error creating session. Please try again.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering session error: %v", err)
		}
		return
	}

	log.Printf("User successfully logged in: %s (ID: %s)", authResponse.Record.Email, authResponse.Record.Id)

	// Redirect to dashboard
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func AuthRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	// Basic validation
	if name == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("All fields are required.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering validation error: %v", err)
		}
		return
	}

	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("Passwords do not match.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering validation error: %v", err)
		}
		return
	}

	// Register user with PocketBase
	user, err := auth.RegisterUser(name, email, password, confirmPassword, Repo)
	if err != nil {
		log.Printf("Registration failed for %s: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("Registration failed. Please try again or use a different email.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering registration error: %v", err)
		}
		return
	}

	log.Printf("User successfully registered: %s (ID: %s)", user.Email, user.Id)

	// Show success message and redirect to login
	w.WriteHeader(http.StatusOK)
	err = templates.AuthRegistrationSuccess().Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering registration success: %v", err)
	}
}

func AuthForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")

	// Basic validation
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("Email is required.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering validation error: %v", err)
		}
		return
	}

	// Request password reset from PocketBase
	err = auth.RequestPasswordReset(email)
	if err != nil {
		log.Printf("Password reset failed for %s: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("Failed to send password reset email. Please check your email address.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering password reset error: %v", err)
		}
		return
	}

	log.Printf("Password reset email sent to: %s", email)

	// Render success message for HTMX target
	w.WriteHeader(http.StatusOK)
	err = templates.ForgotPasswordSuccess(email).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func AuthResendResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")

	// Resend password reset email via PocketBase
	err = auth.RequestPasswordReset(email)
	if err != nil {
		log.Printf("Password reset resend failed for %s: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		err := templates.AuthValidationError("Failed to resend reset email. Please try again.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering resend error: %v", err)
		}
		return
	}

	log.Printf("Password reset email resent to: %s", email)

	w.WriteHeader(http.StatusOK)
	err = templates.AuthPasswordResetSuccess().Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering password reset success: %v", err)
	}
}

func AuthLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear the user session
	err := auth.ClearSession(r, w)
	if err != nil {
		log.Printf("Error clearing user session: %v", err)
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	log.Printf("User successfully logged out")

	// Redirect to login page
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func ConfirmEmailChangeHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Store the token in session for later use
	err := auth.SetEmailChangeToken(r, w, token)
	if err != nil {
		log.Printf("Error storing email change token in session: %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	// Render the email change confirmation page
	err = templates.ConfirmEmailChange(token).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func ConfirmEmailChangeWithPasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")

	// Basic validation
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.EmailChangeConfirmationError("Password is required.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering error template: %v", err)
		}
		return
	}

	// Get token from session
	token, tokenExists := auth.GetEmailChangeToken(r)
	if !tokenExists {
		w.WriteHeader(http.StatusBadRequest)
		err := templates.EmailChangeConfirmationError("No email change request found. Please try again from the email link.").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering error template: %v", err)
		}
		return
	}

	// Confirm email change via auth package
	err := auth.ConfirmEmailChange(token, password)
	if err != nil {
		log.Printf("Email change confirmation failed: %v", err)

		if err.Error() == "validation_error" {
			// PocketBase validation error - show error popup
			w.WriteHeader(http.StatusBadRequest)
			err := templates.EmailChangeErrorPopup().Render(r.Context(), w)
			if err != nil {
				log.Printf("Error rendering error popup template: %v", err)
			}
		} else {
			// Other errors (server config, network, etc.)
			w.WriteHeader(http.StatusInternalServerError)
			err := templates.EmailChangeConfirmationError("An error occurred while processing your request.").Render(r.Context(), w)
			if err != nil {
				log.Printf("Error rendering error template: %v", err)
			}
		}
		return
	}

	// Success - clear token from session
	err = auth.ClearEmailChangeToken(r, w)
	if err != nil {
		log.Printf("Error clearing email change token from session: %v", err)
	}

	// Show success popup and redirect to login
	w.WriteHeader(http.StatusOK)
	err = templates.EmailChangeSuccessPopup().Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering success template: %v", err)
	}
}
