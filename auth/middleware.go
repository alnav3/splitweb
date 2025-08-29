package auth

import (
	"log"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, userId, sessionExists := GetUserFromSession(r)
		if !sessionExists {
			redirectToLogin(w, r)
			return
		}

		// Validate the token with PocketBase
		err := ValidateToken(token, userId)
		if err != nil {
			// Token is invalid, clear the session and redirect
			log.Printf("MIDDLEWARE: Token validation failed for user %s: %v", userId, err)
			ClearSession(r, w)
			redirectToLogin(w, r)
			return
		}

		log.Printf("MIDDLEWARE: Token validation successful for user %s", userId)

		next(w, r)
	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	// For HTMX requests, redirect via header
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// For regular requests, redirect normally
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RedirectIfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, userId, sessionExists := GetUserFromSession(r)
		if sessionExists {
			// Validate the token with PocketBase
			err := ValidateToken(token, userId)
			if err == nil {
				// Token is valid, redirect to dashboard
				log.Printf("MIDDLEWARE: User %s already authenticated, redirecting to dashboard", userId)
				redirectToDashboard(w, r)
				return
			} else {
				// Token is invalid, clear the session
				log.Printf("MIDDLEWARE: Invalid token for user %s, clearing session: %v", userId, err)
				ClearSession(r, w)
			}
		}

		next(w, r)
	}
}

func redirectToDashboard(w http.ResponseWriter, r *http.Request) {
	// For HTMX requests, redirect via header
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	// For regular requests, redirect normally
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
