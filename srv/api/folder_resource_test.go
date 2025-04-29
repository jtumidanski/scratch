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

func TestFolderResource_CRUD(t *testing.T) {
	// Setup test database
	db := database.NewTestDB(t)
	defer database.CleanupTestDB(t, db)

	// Create resource
	resource := NewFolderResource(db)

	// Create a test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	require.NoError(t, db.Create(&user).Error, "Failed to create test user")

	// Test Create
	t.Run("Create", func(t *testing.T) {
		folder := models.Folder{
			Name:   "Test Folder",
			UserID: user.ID,
		}

		resp, err := resource.Create(folder, api2go.Request{})
		require.NoError(t, err, "Failed to create folder")
		require.Equal(t, http.StatusCreated, resp.StatusCode(), "Expected status code 201")

		createdFolder, ok := resp.Result().(models.Folder)
		require.True(t, ok, "Expected result to be a Folder")
		assert.NotEqual(t, uuid.Nil, createdFolder.ID, "Expected folder ID to be set")
		assert.Equal(t, "Test Folder", createdFolder.Name, "Expected folder name to match")
		assert.Equal(t, user.ID, createdFolder.UserID, "Expected folder user ID to match")
	})

	// Create a folder for subsequent tests
	folder := models.Folder{
		Name:   "Test Folder",
		UserID: user.ID,
	}
	require.NoError(t, db.Create(&folder).Error, "Failed to create test folder")

	// Test FindOne
	t.Run("FindOne", func(t *testing.T) {
		resp, err := resource.FindOne(folder.ID.String(), api2go.Request{})
		require.NoError(t, err, "Failed to find folder")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		foundFolder, ok := resp.Result().(models.Folder)
		require.True(t, ok, "Expected result to be a Folder")
		assert.Equal(t, folder.ID, foundFolder.ID, "Expected folder ID to match")
		assert.Equal(t, folder.Name, foundFolder.Name, "Expected folder name to match")
		assert.Equal(t, folder.UserID, foundFolder.UserID, "Expected folder user ID to match")
	})

	// Test FindAll
	t.Run("FindAll", func(t *testing.T) {
		resp, err := resource.FindAll(api2go.Request{})
		require.NoError(t, err, "Failed to find all folders")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		folders, ok := resp.Result().([]models.Folder)
		require.True(t, ok, "Expected result to be a slice of Folders")
		assert.GreaterOrEqual(t, len(folders), 1, "Expected at least one folder")

		// Find our test folder
		var found bool
		for _, f := range folders {
			if f.ID == folder.ID {
				found = true
				assert.Equal(t, folder.Name, f.Name, "Expected folder name to match")
				assert.Equal(t, folder.UserID, f.UserID, "Expected folder user ID to match")
				break
			}
		}
		assert.True(t, found, "Expected to find test folder in results")
	})

	// Test nested folders
	t.Run("NestedFolders", func(t *testing.T) {
		// Create a child folder
		childFolder := models.Folder{
			Name:     "Child Folder",
			UserID:   user.ID,
			ParentID: &folder.ID,
		}
		require.NoError(t, db.Create(&childFolder).Error, "Failed to create child folder")

		// Test FindAll with parent_id filter
		req := api2go.Request{
			QueryParams: map[string][]string{
				"parent_id": {folder.ID.String()},
			},
		}
		resp, err := resource.FindAll(req)
		require.NoError(t, err, "Failed to find folders with parent_id filter")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		folders, ok := resp.Result().([]models.Folder)
		require.True(t, ok, "Expected result to be a slice of Folders")
		assert.Equal(t, 1, len(folders), "Expected exactly one folder")
		assert.Equal(t, childFolder.ID, folders[0].ID, "Expected child folder ID to match")
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updatedFolder := models.Folder{
			ID:     folder.ID,
			Name:   "Updated Folder",
			UserID: user.ID, // This will be preserved by the Update method
		}

		resp, err := resource.Update(updatedFolder, api2go.Request{})
		require.NoError(t, err, "Failed to update folder")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		result, ok := resp.Result().(models.Folder)
		require.True(t, ok, "Expected result to be a Folder")
		assert.Equal(t, folder.ID, result.ID, "Expected folder ID to match")
		assert.Equal(t, "Updated Folder", result.Name, "Expected folder name to be updated")
		assert.Equal(t, user.ID, result.UserID, "Expected folder user ID to be preserved")

		// Verify in database
		var dbFolder models.Folder
		require.NoError(t, db.First(&dbFolder, "id = ?", folder.ID).Error, "Failed to find folder in database")
		assert.Equal(t, "Updated Folder", dbFolder.Name, "Expected folder name to be updated in database")
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		resp, err := resource.Delete(folder.ID.String(), api2go.Request{})
		require.NoError(t, err, "Failed to delete folder")
		require.Equal(t, http.StatusNoContent, resp.StatusCode(), "Expected status code 204")

		// Verify folder is deleted
		var count int64
		db.Model(&models.Folder{}).Where("id = ?", folder.ID).Count(&count)
		assert.Equal(t, int64(0), count, "Expected folder to be deleted")
	})
}