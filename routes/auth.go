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
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Email and password are required.</p>
		</div>`))
		return
	}

	// Authenticate with PocketBase
	authResponse, err := auth.AuthWithPassword(email, password, Repo)
	if err != nil {
		log.Printf("Authentication failed for %s: %v", email, err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Invalid email or password.</p>
		</div>`))
		return
	}

	// Set user session
	err = auth.SetUserSession(r, w, authResponse.Token, authResponse.Record.Id)
	if err != nil {
		log.Printf("Error setting user session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Error creating session. Please try again.</p>
		</div>`))
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
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">All fields are required.</p>
		</div>`))
		return
	}

	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Passwords do not match.</p>
		</div>`))
		return
	}

	// Register user with PocketBase
	user, err := auth.RegisterUser(name, email, password, confirmPassword, Repo)
	if err != nil {
		log.Printf("Registration failed for %s: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Registration failed. Please try again or use a different email.</p>
		</div>`))
		return
	}

	log.Printf("User successfully registered: %s (ID: %s)", user.Email, user.Id)

	// Show success message and redirect to login
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 rounded-lg p-4 mb-4">
		<p class="text-sm text-green-700 dark:text-green-300">Registration successful! Redirecting to login...</p>
	</div>
	<script>setTimeout(function() { window.location.href = '/login'; }, 2000);</script>`))
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
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Email is required.</p>
		</div>`))
		return
	}

	// Request password reset from PocketBase
	err = auth.RequestPasswordReset(email)
	if err != nil {
		log.Printf("Password reset failed for %s: %v", email, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4 mb-4">
			<p class="text-sm text-red-700 dark:text-red-300">Failed to send password reset email. Please check your email address.</p>
		</div>`))
		return
	}

	log.Printf("Password reset email sent to: %s", email)

	// Render success message
	err = templates.PasswordResetSent(email).Render(r.Context(), w)
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
		w.Write([]byte(`<div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-4">
			<p class="text-sm text-red-700 dark:text-red-300">Failed to resend reset email. Please try again.</p>
		</div>`))
		return
	}

	log.Printf("Password reset email resent to: %s", email)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 rounded-lg p-4">
		<p class="text-sm text-green-700 dark:text-green-300">Reset email sent successfully!</p>
	</div>`))
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
