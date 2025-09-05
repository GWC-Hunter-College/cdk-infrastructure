package shared

type EventStatus string

const (
	StatusDrafted  EventStatus = "drafted"
	StatusPosted   EventStatus = "posted"
	StatusArchived EventStatus = "archived"
)

type Event struct {
	ID       int         `json:"id"`
	Title    string      `json:"title"`
	Location string      `json:"location,omitempty"`
	RSVPLink string      `json:"rsvpLink,omitempty"`
	Status   EventStatus `json:"status"`
	StartISO string      `json:"startDate"`
	EndISO   string      `json:"endDate"`
	// Timezone string      `json:"timezone"`
}

type Events struct {
	Events []Event `json:"clubs"`
}

type Club struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Clubs struct {
	Clubs []Club `json:"clubs"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
