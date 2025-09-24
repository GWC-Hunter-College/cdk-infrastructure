package models

type Event struct {
	EventID     int     `json:"id" db:"id"`
	AuthorID    string  `json:"authorId" db:"fk_author_id"`
	ThumbnailID *string `json:"thumbnailId" db:"fk_thumbnail_id"`
	Title       string  `json:"title" db:"title"`
	Location    string  `json:"location" db:"location"`
	RsvpLink    string  `json:"rsvpLink" db:"rsvp_link"`
	Status      string  `json:"status" db:"status"`
	StartDate   string  `json:"startDate" db:"start_date"`
	EndDate     string  `json:"endDate" db:"end_date"`
	Timezone    string  `json:"timezone" db:"timezone"`
	CreatedAt   string  `json:"createdAt" db:"created_at"`
	UpdatedAt   string  `json:"updatedAt" db:"updated_at"`
	DeletedAt   *string `json:"deletedAt" db:"deleted_at"`
}

type EventDescription struct {
	EventID     int     `json:"eventId" db:"fk_event_id"`
	Description *string `json:"description,omitempty" db:"description"`
}
