package api

import (
	"net/http"
	"srv/database"
	"srv/models"
	"testing"

	"github.com/google/uuid"
	"github.com/manyminds/api2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResource_CRUD(t *testing.T) {
	// Setup test database
	db := database.NewTestDB(t)
	defer database.CleanupTestDB(t, db)

	// Create resource
	resource := NewUserResource(db)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Email:    "test@example.com",
		}

		resp, err := resource.Create(user, api2go.Request{})
		require.NoError(t, err, "Failed to create user")
		require.Equal(t, http.StatusCreated, resp.StatusCode(), "Expected status code 201")

		createdUser, ok := resp.Result().(models.User)
		require.True(t, ok, "Expected result to be a User")
		assert.NotEqual(t, uuid.Nil, createdUser.ID, "Expected user ID to be set")
		assert.Equal(t, "testuser", createdUser.Username, "Expected username to match")
		assert.Equal(t, "test@example.com", createdUser.Email, "Expected email to match")
	})

	// Create a user for subsequent tests
	user := models.User{
		Username: "testuser2",
		Email:    "test2@example.com",
	}
	require.NoError(t, db.Create(&user).Error, "Failed to create test user")

	// Test FindOne
	t.Run("FindOne", func(t *testing.T) {
		resp, err := resource.FindOne(user.ID.String(), api2go.Request{})
		require.NoError(t, err, "Failed to find user")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		foundUser, ok := resp.Result().(models.User)
		require.True(t, ok, "Expected result to be a User")
		assert.Equal(t, user.ID, foundUser.ID, "Expected user ID to match")
		assert.Equal(t, user.Username, foundUser.Username, "Expected username to match")
		assert.Equal(t, user.Email, foundUser.Email, "Expected email to match")
	})

	// Test FindAll
	t.Run("FindAll", func(t *testing.T) {
		resp, err := resource.FindAll(api2go.Request{})
		require.NoError(t, err, "Failed to find all users")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		users, ok := resp.Result().([]models.User)
		require.True(t, ok, "Expected result to be a slice of Users")
		assert.GreaterOrEqual(t, len(users), 1, "Expected at least one user")

		// Find our test user
		var found bool
		for _, u := range users {
			if u.ID == user.ID {
				found = true
				assert.Equal(t, user.Username, u.Username, "Expected username to match")
				assert.Equal(t, user.Email, u.Email, "Expected email to match")
				break
			}
		}
		assert.True(t, found, "Expected to find test user in results")
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updatedUser := models.User{
			ID:       user.ID,
			Username: "updateduser",
			Email:    "updated@example.com",
		}

		resp, err := resource.Update(updatedUser, api2go.Request{})
		require.NoError(t, err, "Failed to update user")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		result, ok := resp.Result().(models.User)
		require.True(t, ok, "Expected result to be a User")
		assert.Equal(t, user.ID, result.ID, "Expected user ID to match")
		assert.Equal(t, "updateduser", result.Username, "Expected username to be updated")
		assert.Equal(t, "updated@example.com", result.Email, "Expected email to be updated")

		// Verify in database
		var dbUser models.User
		require.NoError(t, db.First(&dbUser, "id = ?", user.ID).Error, "Failed to find user in database")
		assert.Equal(t, "updateduser", dbUser.Username, "Expected username to be updated in database")
		assert.Equal(t, "updated@example.com", dbUser.Email, "Expected email to be updated in database")
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		resp, err := resource.Delete(user.ID.String(), api2go.Request{})
		require.NoError(t, err, "Failed to delete user")
		require.Equal(t, http.StatusNoContent, resp.StatusCode(), "Expected status code 204")

		// Verify user is deleted
		var count int64
		db.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
		assert.Equal(t, int64(0), count, "Expected user to be deleted")
	})

	// Test unique constraints
	t.Run("UniqueConstraints", func(t *testing.T) {
		// Create a user
		user1 := models.User{
			Username: "uniqueuser",
			Email:    "unique@example.com",
		}
		resp, err := resource.Create(user1, api2go.Request{})
		require.NoError(t, err, "Failed to create first user")
		require.Equal(t, http.StatusCreated, resp.StatusCode(), "Expected status code 201")

		// Try to create another user with the same username
		user2 := models.User{
			Username: "uniqueuser",
			Email:    "different@example.com",
		}
		_, err = resource.Create(user2, api2go.Request{})
		require.Error(t, err, "Expected error when creating user with duplicate username")
		_, ok := err.(api2go.HTTPError)
		require.True(t, ok, "Expected error to be an HTTPError")

		// Try to create another user with the same email
		user3 := models.User{
			Username: "differentuser",
			Email:    "unique@example.com",
		}
		_, err = resource.Create(user3, api2go.Request{})
		require.Error(t, err, "Expected error when creating user with duplicate email")
		_, ok = err.(api2go.HTTPError)
		require.True(t, ok, "Expected error to be an HTTPError")
	})
}
