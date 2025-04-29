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

func TestDocumentResource_CRUD(t *testing.T) {
	// Setup test database
	db := database.NewTestDB(t)
	defer database.CleanupTestDB(t, db)

	// Create resource
	resource := NewDocumentResource(db)

	// Create a test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	require.NoError(t, db.Create(&user).Error, "Failed to create test user")

	// Test Create
	t.Run("Create", func(t *testing.T) {
		doc := models.Document{
			Title:   "Test Document",
			Content: "Test Content",
			UserID:  user.ID,
		}

		resp, err := resource.Create(doc, api2go.Request{})
		require.NoError(t, err, "Failed to create document")
		require.Equal(t, http.StatusCreated, resp.StatusCode(), "Expected status code 201")

		createdDoc, ok := resp.Result().(models.Document)
		require.True(t, ok, "Expected result to be a Document")
		assert.NotEqual(t, uuid.Nil, createdDoc.ID, "Expected document ID to be set")
		assert.Equal(t, "Test Document", createdDoc.Title, "Expected document title to match")
		assert.Equal(t, "Test Content", createdDoc.Content, "Expected document content to match")
		assert.Equal(t, user.ID, createdDoc.UserID, "Expected document user ID to match")
	})

	// Create a document for subsequent tests
	doc := models.Document{
		Title:   "Test Document",
		Content: "Test Content",
		UserID:  user.ID,
	}
	require.NoError(t, db.Create(&doc).Error, "Failed to create test document")

	// Test FindOne
	t.Run("FindOne", func(t *testing.T) {
		resp, err := resource.FindOne(doc.ID.String(), api2go.Request{})
		require.NoError(t, err, "Failed to find document")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		foundDoc, ok := resp.Result().(models.Document)
		require.True(t, ok, "Expected result to be a Document")
		assert.Equal(t, doc.ID, foundDoc.ID, "Expected document ID to match")
		assert.Equal(t, doc.Title, foundDoc.Title, "Expected document title to match")
		assert.Equal(t, doc.Content, foundDoc.Content, "Expected document content to match")
		assert.Equal(t, doc.UserID, foundDoc.UserID, "Expected document user ID to match")
	})

	// Test FindAll
	t.Run("FindAll", func(t *testing.T) {
		resp, err := resource.FindAll(api2go.Request{})
		require.NoError(t, err, "Failed to find all documents")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		docs, ok := resp.Result().([]models.Document)
		require.True(t, ok, "Expected result to be a slice of Documents")
		assert.GreaterOrEqual(t, len(docs), 1, "Expected at least one document")

		// Find our test document
		var found bool
		for _, d := range docs {
			if d.ID == doc.ID {
				found = true
				assert.Equal(t, doc.Title, d.Title, "Expected document title to match")
				assert.Equal(t, doc.Content, d.Content, "Expected document content to match")
				assert.Equal(t, doc.UserID, d.UserID, "Expected document user ID to match")
				break
			}
		}
		assert.True(t, found, "Expected to find test document in results")
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updatedDoc := models.Document{
			ID:      doc.ID,
			Title:   "Updated Document",
			Content: "Updated Content",
			UserID:  user.ID, // This will be preserved by the Update method
		}

		resp, err := resource.Update(updatedDoc, api2go.Request{})
		require.NoError(t, err, "Failed to update document")
		require.Equal(t, http.StatusOK, resp.StatusCode(), "Expected status code 200")

		result, ok := resp.Result().(models.Document)
		require.True(t, ok, "Expected result to be a Document")
		assert.Equal(t, doc.ID, result.ID, "Expected document ID to match")
		assert.Equal(t, "Updated Document", result.Title, "Expected document title to be updated")
		assert.Equal(t, "Updated Content", result.Content, "Expected document content to be updated")
		assert.Equal(t, user.ID, result.UserID, "Expected document user ID to be preserved")

		// Verify in database
		var dbDoc models.Document
		require.NoError(t, db.First(&dbDoc, "id = ?", doc.ID).Error, "Failed to find document in database")
		assert.Equal(t, "Updated Document", dbDoc.Title, "Expected document title to be updated in database")
		assert.Equal(t, "Updated Content", dbDoc.Content, "Expected document content to be updated in database")
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		resp, err := resource.Delete(doc.ID.String(), api2go.Request{})
		require.NoError(t, err, "Failed to delete document")
		require.Equal(t, http.StatusNoContent, resp.StatusCode(), "Expected status code 204")

		// Verify document is deleted
		var count int64
		db.Model(&models.Document{}).Where("id = ?", doc.ID).Count(&count)
		assert.Equal(t, int64(0), count, "Expected document to be deleted")
	})
}