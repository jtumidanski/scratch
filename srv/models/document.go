package models

import (
	"github.com/google/uuid"
	"github.com/manyminds/api2go/jsonapi"
	"gorm.io/gorm"
	"time"
)

// Document represents a document in the system
type Document struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	FolderID  *uuid.UUID     `gorm:"type:uuid;null" json:"folder_id"`
	Folder    *Folder        `gorm:"foreignKey:FolderID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (d Document) GetID() string {
	return d.ID.String()
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (d *Document) SetID(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	d.ID = uuid
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (d Document) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "user",
		},
		{
			Type: "folders",
			Name: "folder",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (d Document) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{
		{
			ID:   d.UserID.String(),
			Type: "users",
			Name: "user",
		},
	}

	// Add folder if exists
	if d.FolderID != nil {
		result = append(result, jsonapi.ReferenceID{
			ID:   d.FolderID.String(),
			Type: "folders",
			Name: "folder",
		})
	}

	return result
}

// BeforeCreate will set a UUID rather than numeric ID
func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
