package routes

import (
	"log"
	"net/http"

	"github.com/alnav3/splitweb/templates"
)

func GroupDetailHandler(w http.ResponseWriter, r *http.Request) {

	// Extract group ID from URL
	groupID := r.PathValue("id")

	// Get mock data
	groupDetail := getMockGroupDetail(groupID)

	// Render the template
	err := templates.GroupDetailPage(groupDetail).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

// HTMX tab handlers
func ChargesTabHandler(w http.ResponseWriter, r *http.Request) {

	groupDetail := getMockGroupDetail("1")
	err := templates.ChargesTabContent(groupDetail).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func BalancesTabHandler(w http.ResponseWriter, r *http.Request) {
	groupDetail := getMockGroupDetail("1")
	err := templates.BalancesTabContent(groupDetail).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func ActivityTabHandler(w http.ResponseWriter, r *http.Request) {
	groupDetail := getMockGroupDetail("1")
	err := templates.ActivityTabContent(groupDetail).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func BreakdownTabHandler(w http.ResponseWriter, r *http.Request) {
	groupDetail := getMockGroupDetail("1")
	err := templates.BreakdownTabContent(groupDetail).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

// HTMX form handlers
func AddChargeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// In a real app, we would save the charge to the database
	// For now, just return a success message or render the new charge
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Charge added successfully"))
}

func AddSettlementHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// In a real app, we would save the settlement to the database
	// For now, just return a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Settlement recorded successfully"))
}

func InviteMemberHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	groupID := r.FormValue("groupId")
	email := r.FormValue("email")
	message := r.FormValue("message")

	// In a real app, we would:
	// 1. Validate the group exists and user has permission to invite
	// 2. Check if email is already a member
	// 3. Generate invitation token and store it
	// 4. Send invitation email
	log.Printf("Group invitation: Group ID: %s, Email: %s, Message: %s", groupID, email, message)

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Invitation sent successfully"))
}

// Helper function to create mock group data
func getMockGroupDetail(groupID string) templates.GroupDetail {

	// Create mock members
	members := []templates.Member{
		{ID: "1", Name: "Alice", ImageURL: "https://ui-avatars.com/api/?name=Alice&background=6366F1&color=fff"},
		{ID: "2", Name: "Bob", ImageURL: "https://ui-avatars.com/api/?name=Bob&background=8B5CF6&color=fff"},
		{ID: "3", Name: "Carol", ImageURL: "https://ui-avatars.com/api/?name=Carol&background=EC4899&color=fff"},
		{ID: "4", Name: "Dave", ImageURL: "https://ui-avatars.com/api/?name=Dave&background=F59E0B&color=fff"},
	}

	// Create mock charges
	charges := []templates.Charge{
		{
			ID:           "1",
			Description:  "Dinner at Italian Restaurant",
			Amount:       120.00,
			PaidBy:       members[0],
			Date:         "2023-08-15",
			Participants: members,
		},
		{
			ID:           "2",
			Description:  "Grocery Shopping",
			Amount:       85.75,
			PaidBy:       members[1],
			Date:         "2023-08-10",
			Participants: []templates.Member{members[0], members[1], members[2]},
		},
		{
			ID:           "3",
			Description:  "Movie Tickets",
			Amount:       48.00,
			PaidBy:       members[2],
			Date:         "2023-08-05",
			Participants: []templates.Member{members[0], members[2], members[3]},
		},
	}

	// Create mock debts
	debts := []templates.Debt{
		{FromMember: members[0], ToMember: members[1], Amount: 25.50},
		{FromMember: members[0], ToMember: members[2], Amount: 12.75},
		{FromMember: members[1], ToMember: members[3], Amount: 30.00},
		{FromMember: members[3], ToMember: members[0], Amount: 18.25},
	}

	// Create detailed debts for the matrix
	allDebts := []templates.Debt{}
	for _, from := range members {
		for _, to := range members {
			if from.ID != to.ID {
				amount := 0.0
				for _, debt := range debts {
					if debt.FromMember.ID == from.ID && debt.ToMember.ID == to.ID {
						amount = debt.Amount
						break
					} else if debt.FromMember.ID == to.ID && debt.ToMember.ID == from.ID {
						amount = -debt.Amount
						break
					}
				}
				allDebts = append(allDebts, templates.Debt{FromMember: from, ToMember: to, Amount: amount})
			}
		}
	}

	return templates.GroupDetail{
		ID:          groupID,
		Name:        "Summer Trip 2023",
		Description: "Trip to the beach with friends",
		ImageURL:    "https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?ixlib=rb-1.2.1&auto=format&fit=crop&w=1350&q=80",
		Members:     members,
		Charges:     charges,
		Payments:    []templates.Payment{},
		Debts:       allDebts,
	}
}
