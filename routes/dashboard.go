package routes

import (
	"log"
	"net/http"

	"github.com/alnav3/splitweb/templates"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	err := templates.Dashboard().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}
