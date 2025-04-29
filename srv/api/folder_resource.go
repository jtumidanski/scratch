package api

import (
	"net/http"
	"srv/models"

	"github.com/google/uuid"
	"github.com/manyminds/api2go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FolderResource implements api2go.CRUD interface for Folder
type FolderResource struct {
	DB *gorm.DB
}

// NewFolderResource creates a new FolderResource
func NewFolderResource(db *gorm.DB) *FolderResource {
	return &FolderResource{
		DB: db,
	}
}

// FindAll returns all folders
func (r FolderResource) FindAll(req api2go.Request) (api2go.Responder, error) {
	logrus.Info("Finding all folders")

	var folders []models.Folder
	query := r.DB

	// Filter by user ID if provided
	if userID, ok := req.QueryParams["user_id"]; ok && len(userID) > 0 {
		logrus.WithField("user_id", userID[0]).Info("Filtering folders by user ID")

		uuid, err := uuid.Parse(userID[0])
		if err != nil {
			logrus.WithError(err).WithField("user_id", userID[0]).Error("Invalid user ID")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid user ID", http.StatusBadRequest)
		}

		query = query.Where("user_id = ?", uuid)
	}

	// Filter by parent ID if provided
	if parentID, ok := req.QueryParams["parent_id"]; ok && len(parentID) > 0 {
		if parentID[0] == "null" {
			// Get root folders (no parent)
			logrus.Info("Filtering folders with no parent")
			query = query.Where("parent_id IS NULL")
		} else {
			logrus.WithField("parent_id", parentID[0]).Info("Filtering folders by parent ID")

			uuid, err := uuid.Parse(parentID[0])
			if err != nil {
				logrus.WithError(err).WithField("parent_id", parentID[0]).Error("Invalid parent ID")
				return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid parent ID", http.StatusBadRequest)
			}

			query = query.Where("parent_id = ?", uuid)
		}
	}

	if err := query.Find(&folders).Error; err != nil {
		logrus.WithError(err).Error("Failed to find folders")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: folders}, nil
}

// FindOne returns a single folder
func (r FolderResource) FindOne(id string, req api2go.Request) (api2go.Responder, error) {
	logrus.WithField("id", id).Info("Finding folder")

	uuid, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Invalid folder ID")
		return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid folder ID", http.StatusBadRequest)
	}

	var folder models.Folder
	if err := r.DB.First(&folder, "id = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", id).Warn("Folder not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Folder not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", id).Error("Failed to find folder")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: folder}, nil
}

// Create creates a new folder
func (r FolderResource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	folder, ok := obj.(models.Folder)
	if !ok {
		err := api2go.NewHTTPError(nil, "Invalid instance given", http.StatusBadRequest)
		logrus.WithError(err).Error("Invalid instance given to create folder")
		return &api2go.Response{}, err
	}

	logrus.WithFields(logrus.Fields{
		"name":      folder.Name,
		"user_id":   folder.UserID,
		"parent_id": folder.ParentID,
	}).Info("Creating folder")

	// Validate user exists
	var user models.User
	if err := r.DB.First(&user, "id = ?", folder.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("user_id", folder.UserID).Warn("User not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "User not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("user_id", folder.UserID).Error("Failed to find user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Validate parent folder exists if provided
	if folder.ParentID != nil {
		var parentFolder models.Folder
		if err := r.DB.First(&parentFolder, "id = ?", folder.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				logrus.WithField("parent_id", folder.ParentID).Warn("Parent folder not found")
				return &api2go.Response{}, api2go.NewHTTPError(err, "Parent folder not found", http.StatusNotFound)
			}
			logrus.WithError(err).WithField("parent_id", folder.ParentID).Error("Failed to find parent folder")
			return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
		}
	}

	if err := r.DB.Create(&folder).Error; err != nil {
		logrus.WithError(err).Error("Failed to create folder")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: folder, Code: http.StatusCreated}, nil
}

// Delete deletes a folder
func (r FolderResource) Delete(id string, req api2go.Request) (api2go.Responder, error) {
	logrus.WithField("id", id).Info("Deleting folder")

	uuid, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Invalid folder ID")
		return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid folder ID", http.StatusBadRequest)
	}

	// Check if folder exists
	var folder models.Folder
	if err := r.DB.First(&folder, "id = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", id).Warn("Folder not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Folder not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", id).Error("Failed to find folder")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Check if folder has subfolders
	var subfolderCount int64
	if err := r.DB.Model(&models.Folder{}).Where("parent_id = ?", uuid).Count(&subfolderCount).Error; err != nil {
		logrus.WithError(err).WithField("id", id).Error("Failed to count subfolders")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	if subfolderCount > 0 {
		err := api2go.NewHTTPError(nil, "Cannot delete folder with subfolders", http.StatusBadRequest)
		logrus.WithField("id", id).Warn("Cannot delete folder with subfolders")
		return &api2go.Response{}, err
	}

	// Check if folder has documents
	var documentCount int64
	if err := r.DB.Model(&models.Document{}).Where("folder_id = ?", uuid).Count(&documentCount).Error; err != nil {
		logrus.WithError(err).WithField("id", id).Error("Failed to count documents")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	if documentCount > 0 {
		err := api2go.NewHTTPError(nil, "Cannot delete folder with documents", http.StatusBadRequest)
		logrus.WithField("id", id).Warn("Cannot delete folder with documents")
		return &api2go.Response{}, err
	}

	// Delete folder
	if err := r.DB.Delete(&folder).Error; err != nil {
		logrus.WithError(err).WithField("id", id).Error("Failed to delete folder")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Code: http.StatusNoContent}, nil
}

// Update updates a folder
func (r FolderResource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	folder, ok := obj.(models.Folder)
	if !ok {
		err := api2go.NewHTTPError(nil, "Invalid instance given", http.StatusBadRequest)
		logrus.WithError(err).Error("Invalid instance given to update folder")
		return &api2go.Response{}, err
	}

	logrus.WithFields(logrus.Fields{
		"id":        folder.ID,
		"name":      folder.Name,
		"parent_id": folder.ParentID,
	}).Info("Updating folder")

	// Check if folder exists
	var existingFolder models.Folder
	if err := r.DB.First(&existingFolder, "id = ?", folder.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", folder.ID).Warn("Folder not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "Folder not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", folder.ID).Error("Failed to find folder")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Validate parent folder exists if provided
	if folder.ParentID != nil {
		// Prevent circular reference
		if *folder.ParentID == folder.ID {
			err := api2go.NewHTTPError(nil, "Folder cannot be its own parent", http.StatusBadRequest)
			logrus.WithField("id", folder.ID).Warn("Folder cannot be its own parent")
			return &api2go.Response{}, err
		}

		var parentFolder models.Folder
		if err := r.DB.First(&parentFolder, "id = ?", folder.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				logrus.WithField("parent_id", folder.ParentID).Warn("Parent folder not found")
				return &api2go.Response{}, api2go.NewHTTPError(err, "Parent folder not found", http.StatusNotFound)
			}
			logrus.WithError(err).WithField("parent_id", folder.ParentID).Error("Failed to find parent folder")
			return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
		}
	}

	// Preserve the user ID
	folder.UserID = existingFolder.UserID

	// Update folder
	if err := r.DB.Save(&folder).Error; err != nil {
		logrus.WithError(err).WithField("id", folder.ID).Error("Failed to update folder")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: folder, Code: http.StatusOK}, nil
}
