package models

import (
	"github.com/google/uuid"
	"github.com/manyminds/api2go/jsonapi"
	"gorm.io/gorm"
	"time"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Username  string         `gorm:"size:255;not null;unique" json:"username"`
	Email     string         `gorm:"size:255;not null;unique" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Folders   []Folder       `gorm:"foreignKey:UserID" json:"-"`
	Documents []Document     `gorm:"foreignKey:UserID" json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u User) GetID() string {
	return u.ID.String()
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *User) SetID(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	u.ID = uuid
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u User) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "folders",
			Name: "folders",
		},
		{
			Type: "documents",
			Name: "documents",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u User) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, folder := range u.Folders {
		result = append(result, jsonapi.ReferenceID{
			ID:   folder.ID.String(),
			Type: "folders",
			Name: "folders",
		})
	}
	for _, document := range u.Documents {
		result = append(result, jsonapi.ReferenceID{
			ID:   document.ID.String(),
			Type: "documents",
			Name: "documents",
		})
	}
	return result
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
