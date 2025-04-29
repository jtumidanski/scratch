package api

import (
	"net/http"
	"srv/models"

	"github.com/google/uuid"
	"github.com/manyminds/api2go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserResource implements api2go.CRUD interface for User
type UserResource struct {
	DB *gorm.DB
}

// NewUserResource creates a new UserResource
func NewUserResource(db *gorm.DB) *UserResource {
	return &UserResource{
		DB: db,
	}
}

// FindAll returns all users
func (r UserResource) FindAll(req api2go.Request) (api2go.Responder, error) {
	logrus.Info("Finding all users")

	var users []models.User
	query := r.DB.Unscoped()

	if err := query.Find(&users).Error; err != nil {
		logrus.WithError(err).Error("Failed to find users")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: users}, nil
}

// FindOne returns a single user
func (r UserResource) FindOne(id string, req api2go.Request) (api2go.Responder, error) {
	logrus.WithField("id", id).Info("Finding user")

	uuid, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Invalid user ID")
		return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid user ID", http.StatusBadRequest)
	}

	var user models.User
	if err := r.DB.First(&user, "id = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", id).Warn("User not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "User not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", id).Error("Failed to find user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: user}, nil
}

// Create creates a new user
func (r UserResource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(models.User)
	if !ok {
		err := api2go.NewHTTPError(nil, "Invalid instance given", http.StatusBadRequest)
		logrus.WithError(err).Error("Invalid instance given to create user")
		return &api2go.Response{}, err
	}

	logrus.WithFields(logrus.Fields{
		"username": user.Username,
		"email":    user.Email,
	}).Info("Creating user")

	if err := r.DB.Create(&user).Error; err != nil {
		logrus.WithError(err).Error("Failed to create user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: user, Code: http.StatusCreated}, nil
}

// Delete deletes a user
func (r UserResource) Delete(id string, req api2go.Request) (api2go.Responder, error) {
	logrus.WithField("id", id).Info("Deleting user")

	uuid, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Invalid user ID")
		return &api2go.Response{}, api2go.NewHTTPError(err, "Invalid user ID", http.StatusBadRequest)
	}

	// Check if user exists
	var user models.User
	if err := r.DB.First(&user, "id = ?", uuid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", id).Warn("User not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "User not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", id).Error("Failed to find user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Delete user
	if err := r.DB.Delete(&user).Error; err != nil {
		logrus.WithError(err).WithField("id", id).Error("Failed to delete user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Code: http.StatusNoContent}, nil
}

// Update updates a user
func (r UserResource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(models.User)
	if !ok {
		err := api2go.NewHTTPError(nil, "Invalid instance given", http.StatusBadRequest)
		logrus.WithError(err).Error("Invalid instance given to update user")
		return &api2go.Response{}, err
	}

	logrus.WithFields(logrus.Fields{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}).Info("Updating user")

	// Check if user exists
	var existingUser models.User
	if err := r.DB.First(&existingUser, "id = ?", user.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("id", user.ID).Warn("User not found")
			return &api2go.Response{}, api2go.NewHTTPError(err, "User not found", http.StatusNotFound)
		}
		logrus.WithError(err).WithField("id", user.ID).Error("Failed to find user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	// Update user
	if err := r.DB.Save(&user).Error; err != nil {
		logrus.WithError(err).WithField("id", user.ID).Error("Failed to update user")
		return &api2go.Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}

	return &api2go.Response{Res: user, Code: http.StatusOK}, nil
}
