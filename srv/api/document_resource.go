package api

import (
	"net/http"
	"srv/models"

	"github.com/google/uuid"
	"github.com/manyminds/api2go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// DocumentResource implements api2go.CRUD interface for Document
type DocumentResource struct {
	DB *gorm.DB
}

// NewDocumentResource creates a new DocumentResource
func NewDocumentResource(db *gorm.DB) *DocumentResource {
	return &DocumentResource{
		DB: db,
	}
}

// FindAll returns all documents
func (r DocumentResource) FindAll(req api2go.Request) (api2go.Responder, error) {
	logrus.Info("Finding all documents")

	var documents []models.Document
	query := r.DB

	// Filter by user ID if provided
	if userID, ok := req.QueryParams["user_id"]; ok && len(userID) > 0 {
		logrus.WithField("user_id", userID[0]).Info("Filtering documents by user ID")

		uuid, err := uuid.Parse(userID[0])
		if err != nil {
			logrus.WithError(err).WithField("user_id", userID[0]).Error("Invalid user ID")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid user ID", http.StatusBadRequest)
		}

		query = query.Where("user_id = ?", uuid)
	}

	// Filter by folder ID if provided
	if folderID, ok := req.QueryParams["folder_id"]; ok && len(folderID) > 0 {
		if folderID[0] == "null" {
			// Get documents with no folder
			logrus.Info("Filtering documents with no folder")
			query = query.Where("folder_id IS NULL")
		} else {
			logrus.WithField("folder_id", folderID[0]).Info("Filtering documents by folder ID")

			uuid, err := uuid.Parse(folderID[0])
			if err != nil {
				logrus.WithError(err).WithField("folder_id", folderID[0]).Error("Invalid folder ID")
				return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid folder ID", http.StatusBadRequest)
			}

			query = query.Where("folder_id = ?", uuid)
		}
	}

	if err := query.Find(&documents).Error; err != nil {
		logrus.WithError(err).Error("Failed to find documents")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: documents}, nil
}

// FindOne returns a single document
func (r DocumentResource) FindOne(id string, req api2go.Request) (api2go.Responder, error) {
	logrus.WithField("id", id).Info("Finding document")

	uuid, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Invalid document ID")
		return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid document ID", http.StatusBadRequest)
	}

	var document models.Document
	if err := r.DB.First(&document, "id = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", id).Warn("Document not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Document not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", id).Error("Failed to find document")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: document}, nil
}

// Create creates a new document
func (r DocumentResource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	document, ok := obj.(models.Document)
	if !ok {
		err := api2go.NewHTTPError(nil, "Invalid instance given", http.StatusBadRequest)
		logrus.WithError(err).Error("Invalid instance given to create document")
		return &api2go.Response{}, err
	}

	logrus.WithFields(logrus.Fields{
		"title":     document.Title,
		"user_id":   document.UserID,
		"folder_id": document.FolderID,
	}).Info("Creating document")

	// Validate user exists
	var user models.User
	if err := r.DB.First(&user, "id = ?", document.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("user_id", document.UserID).Warn("User not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "User not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("user_id", document.UserID).Error("Failed to find user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Validate folder exists if provided
	if document.FolderID != nil {
		var folder models.Folder
		if err := r.DB.First(&folder, "id = ?", document.FolderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				logrus.WithField("folder_id", document.FolderID).Warn("Folder not found")
				return &api2go.Response{}, api2go.NewHTTPError(err, "Folder not found", http.StatusNotFound)
			}
			logrus.WithError(err).WithField("folder_id", document.FolderID).Error("Failed to find folder")
			return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
		}

		// Ensure folder belongs to the same user
		if folder.UserID != document.UserID {
			err := api2go.NewHTTPError(nil, "Folder does not belong to the user", http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"folder_id": document.FolderID,
				"user_id":   document.UserID,
			}).Warn("Folder does not belong to the user")
			return &api2go.Response{}, err
		}
	}

	if err := r.DB.Create(&document).Error; err != nil {
		logrus.WithError(err).Error("Failed to create document")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: document, Code: http.StatusCreated}, nil
}

// Delete deletes a document
func (r DocumentResource) Delete(id string, req api2go.Request) (api2go.Responder, error) {
	logrus.WithField("id", id).Info("Deleting document")

	uuid, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Invalid document ID")
		return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid document ID", http.StatusBadRequest)
	}

	// Check if document exists
	var document models.Document
	if err := r.DB.First(&document, "id = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", id).Warn("Document not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Document not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", id).Error("Failed to find document")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Delete document
	if err := r.DB.Delete(&document).Error; err != nil {
		logrus.WithError(err).WithField("id", id).Error("Failed to delete document")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Code: http.StatusNoContent}, nil
}

// Update updates a document
func (r DocumentResource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	document, ok := obj.(models.Document)
	if !ok {
		err := api2go.NewHTTPError(nil, "Invalid instance given", http.StatusBadRequest)
		logrus.WithError(err).Error("Invalid instance given to update document")
		return &api2go.Response{}, err
	}

	logrus.WithFields(logrus.Fields{
		"id":        document.ID,
		"title":     document.Title,
		"folder_id": document.FolderID,
	}).Info("Updating document")

	// Check if document exists
	var existingDocument models.Document
	if err := r.DB.First(&existingDocument, "id = ?", document.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", document.ID).Warn("Document not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Document not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", document.ID).Error("Failed to find document")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Validate folder exists if provided
	if document.FolderID != nil {
		var folder models.Folder
		if err := r.DB.First(&folder, "id = ?", document.FolderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				logrus.WithField("folder_id", document.FolderID).Warn("Folder not found")
				return &api2go.Response{}, api2go.NewHTTPError(err, "Folder not found", http.StatusNotFound)
			}
			logrus.WithError(err).WithField("folder_id", document.FolderID).Error("Failed to find folder")
			return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
		}

		// Ensure folder belongs to the same user
		if folder.UserID != existingDocument.UserID {
			err := api2go.NewHTTPError(nil, "Folder does not belong to the user", http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"folder_id": document.FolderID,
				"user_id":   existingDocument.UserID,
			}).Warn("Folder does not belong to the user")
			return &api2go.Response{}, err
		}
	}

	// Preserve the user ID
	document.UserID = existingDocument.UserID

	// Update document
	if err := r.DB.Save(&document).Error; err != nil {
		logrus.WithError(err).WithField("id", document.ID).Error("Failed to update document")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: document, Code: http.StatusOK}, nil
}
