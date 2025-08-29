package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}
	store = sessions.NewCookieStore([]byte(secretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "splitweb-session")
}

func SaveSession(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return session.Save(r, w)
}

func SetUserSession(r *http.Request, w http.ResponseWriter, token string, userId string) error {
	session, err := GetSession(r)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	session.Values["token"] = token
	session.Values["user_id"] = userId
	session.Values["authenticated"] = true

	return SaveSession(r, w, session)
}

func ClearSession(r *http.Request, w http.ResponseWriter) error {
	session, err := GetSession(r)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	session.Values["token"] = nil
	session.Values["user_id"] = nil
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	return SaveSession(r, w, session)
}

func GetUserFromSession(r *http.Request) (string, string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", "", false
	}

	token, ok := session.Values["token"].(string)
	if !ok {
		return "", "", false
	}

	userId, ok := session.Values["user_id"].(string)
	if !ok {
		return "", "", false
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		return "", "", false
	}

	return token, userId, true
}

func IsAuthenticated(r *http.Request) bool {
	token, userId, sessionExists := GetUserFromSession(r)
	if !sessionExists {
		return false
	}

	// Validate the token with PocketBase
	err := ValidateToken(token, userId)
	return err == nil
}
