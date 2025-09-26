package models

type Club struct {
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name"`
	ThumbnailURL *string `json:"thumbnailUrl,omitempty" db:"thumbnail_url"`
}

type ClubRole string

const (
	RoleMember ClubRole = "member"
	RoleEboard ClubRole = "eboard"
	RoleOwner  ClubRole = "owner"
)
