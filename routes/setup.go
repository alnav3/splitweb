package routes

import (
	"net/http"
	"strings"

	"github.com/alnav3/splitweb/auth"
	databaselogic "github.com/alnav3/splitweb/db/database_logic"
)

var Repo *databaselogic.Repository

// Custom file server that sets proper MIME types
func StaticHandler(w http.ResponseWriter, r *http.Request) {

	// Remove the /static/ prefix
	path := strings.TrimPrefix(r.URL.Path, "/style/")

	// Serve the file
	http.ServeFile(w, r, path)
}

func HandleDirectory(dirs ...string) {
	for _, dir := range dirs {
		http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
	}
}

// SetupRoutes configures all application routes
func SetupRoutes(repo *databaselogic.Repository) {
	// Store the repository globally for use in route handlers
	Repo = repo

	// Static file handling
	HandleDirectory("/style/")

	// Set up front directory for other static assets
	frontFS := http.FileServer(http.Dir("front"))
	http.Handle("/front/", http.StripPrefix("/front/", frontFS))

	// Main page routes (protected)
	http.HandleFunc("/", auth.RequireAuth(DashboardHandler))
	http.HandleFunc("/group/{id}", auth.RequireAuth(GroupDetailHandler))
	http.HandleFunc("/profile", auth.RequireAuth(ProfileHandler))

	// Authentication page routes (redirect if already authenticated)
	http.HandleFunc("/login", auth.RedirectIfAuthenticated(LoginHandler))
	http.HandleFunc("/register", auth.RedirectIfAuthenticated(RegisterHandler))
	http.HandleFunc("/forgot-password", auth.RedirectIfAuthenticated(ForgotPasswordHandler))

	// HTMX routes for tab navigation (protected)
	http.HandleFunc("/group/tab/charges", auth.RequireAuth(ChargesTabHandler))
	http.HandleFunc("/group/tab/balances", auth.RequireAuth(BalancesTabHandler))
	http.HandleFunc("/group/tab/activity", auth.RequireAuth(ActivityTabHandler))
	http.HandleFunc("/group/tab/breakdown", auth.RequireAuth(BreakdownTabHandler))


	// HTMX routes for forms (protected)
	http.HandleFunc("/group/charges", auth.RequireAuth(AddChargeHandler))
	http.HandleFunc("/group/settlements", auth.RequireAuth(AddSettlementHandler))
	http.HandleFunc("/group/invite", auth.RequireAuth(InviteMemberHandler))

	// Authentication form handlers
	http.HandleFunc("POST /auth/login", AuthLoginHandler)
	http.HandleFunc("/auth/register", AuthRegisterHandler)
	http.HandleFunc("/auth/forgot-password", AuthForgotPasswordHandler)
	http.HandleFunc("/auth/resend-reset", AuthResendResetHandler)
	http.HandleFunc("/auth/logout", AuthLogoutHandler)

	// Profile form handlers (protected)
	http.HandleFunc("/profile/email", auth.RequireAuth(ProfileChangeEmailHandler))
	http.HandleFunc("/profile/password", auth.RequireAuth(ProfileChangePasswordHandler))
	http.HandleFunc("/profile/delete", auth.RequireAuth(ProfileDeleteAccountHandler))
}
