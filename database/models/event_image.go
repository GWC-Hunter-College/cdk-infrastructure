package models

type EventImage struct {
	EventID string `db:"fk_event_id"`
	ImageID string `db:"fk_image_id"`
}
