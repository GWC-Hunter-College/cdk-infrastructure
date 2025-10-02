package models

type Club struct {
	ClubID       string `json:"id"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type EventOwners struct {
	Owner      Club   `json:"owner"`
	Associates []Club `json:"associates"`
}
