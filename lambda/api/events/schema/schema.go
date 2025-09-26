package event_schema

import "cdk-infrastructure/database/models"

type SQLSchema struct {
	ClubID    string  `db:"club_id"`
	IsOwner   bool    `db:"is_owner"`
	ObjectKey *string `db:"object_key"`
	models.Event
	Description *string `db:"description"`
}

type Club struct {
	ClubID       string `json:"id"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type OwnerSchema struct {
	Owner      Club   `json:"owner"`
	Associates []Club `json:"associates"`
}

type ResponseSchema struct {
	models.Event
	Description *string `json:"description"`
	OwnerSchema
}
