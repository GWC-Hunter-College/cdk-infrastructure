package shared

// Exported (capitalized) so other packages can use them.

type EventStatus string

const (
	StatusDrafted  EventStatus = "drafted"
	StatusPosted   EventStatus = "posted"
	StatusArchived EventStatus = "archived"
)

type Club struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Verified bool   `json:"verified,omitempty"`
	Role     string `json:"role,omitempty"` // optional: "member"|"eboard"|"owner"
}

type MeClubsResponse struct {
	Clubs []Club `json:"clubs"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
