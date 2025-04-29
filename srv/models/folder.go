package models

import (
	"github.com/google/uuid"
	"github.com/manyminds/api2go/jsonapi"
	"gorm.io/gorm"
	"time"
)

// Folder represents a folder in the system
type Folder struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;null" json:"parent_id"`
	Parent    *Folder        `gorm:"foreignKey:ParentID" json:"-"`
	Folders   []Folder       `gorm:"foreignKey:ParentID" json:"-"`
	Documents []Document     `gorm:"foreignKey:FolderID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (f Folder) GetID() string {
	return f.ID.String()
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (f *Folder) SetID(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	f.ID = uuid
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (f Folder) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "user",
		},
		{
			Type: "folders",
			Name: "parent",
		},
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
func (f Folder) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{
		{
			ID:   f.UserID.String(),
			Type: "users",
			Name: "user",
		},
	}

	// Add parent folder if exists
	if f.ParentID != nil {
		result = append(result, jsonapi.ReferenceID{
			ID:   f.ParentID.String(),
			Type: "folders",
			Name: "parent",
		})
	}

	// Add child folders
	for _, folder := range f.Folders {
		result = append(result, jsonapi.ReferenceID{
			ID:   folder.ID.String(),
			Type: "folders",
			Name: "folders",
		})
	}

	// Add documents
	for _, document := range f.Documents {
		result = append(result, jsonapi.ReferenceID{
			ID:   document.ID.String(),
			Type: "documents",
			Name: "documents",
		})
	}

	return result
}

// BeforeCreate will set a UUID rather than numeric ID
func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
