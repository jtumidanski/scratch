# Document Storage Service

A server-side component for an application that stores and retrieves document data for users. This service provides a RESTful API for managing users, folders, and documents.

## Features

- User management (create, read, update, delete)
- Folder management (create, read, update, delete)
- Document management (create, read, update, delete)
- Hierarchical folder structure
- JSON:API compliant responses

## Technologies Used

- Go (Golang)
- GORM (ORM for Go)
- PostgreSQL (Database)
- api2go (JSON:API implementation)
- Logrus (Logging)
- Bruno (API testing)
- Docker (Containerization)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Docker and Docker Compose (optional)

### Environment Variables

| Variable | Description | Default | Possible Values |
|----------|-------------|---------|----------------|
| DB_HOST | Database host | localhost | Any valid hostname |
| DB_PORT | Database port | 5432 | Any valid port number |
| DB_USER | Database username | postgres | Any valid username |
| DB_PASSWORD | Database password | postgres | Any valid password |
| DB_NAME | Database name | document_storage | Any valid database name |
| DB_SSLMODE | Database SSL mode | disable | disable, require, verify-ca, verify-full |
| PORT | Server port | 8080 | Any valid port number |
| LOG_LEVEL | Logging level | info | trace, debug, info, warn, error, fatal, panic |

### Running with Docker

1. Clone the repository
2. Navigate to the project directory
3. Run the application using Docker Compose:

```bash
docker-compose up -d
```

The API will be available at http://localhost:8080/v1/

### Running Locally

1. Clone the repository
2. Navigate to the project directory
3. Set up the PostgreSQL database
4. Set the required environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=document_storage
export DB_SSLMODE=disable
export PORT=8080
export LOG_LEVEL=info
```

5. Run the application:

```bash
go run main.go
```

The API will be available at http://localhost:8080/v1/

## API Endpoints

The API follows the JSON:API specification (https://jsonapi.org/).

### Users

#### Create a User

- **URL**: `/v1/users`
- **Method**: `POST`
- **Request Body**:
```json
{
  "data": {
    "type": "users",
    "attributes": {
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

#### Get All Users

- **URL**: `/v1/users`
- **Method**: `GET`

#### Get a User

- **URL**: `/v1/users/{id}`
- **Method**: `GET`

#### Update a User

- **URL**: `/v1/users/{id}`
- **Method**: `PATCH`
- **Request Body**:
```json
{
  "data": {
    "type": "users",
    "id": "{id}",
    "attributes": {
      "username": "updateduser",
      "email": "updated@example.com"
    }
  }
}
```

#### Delete a User

- **URL**: `/v1/users/{id}`
- **Method**: `DELETE`

### Folders

#### Create a Folder

- **URL**: `/v1/folders`
- **Method**: `POST`
- **Request Body**:
```json
{
  "data": {
    "type": "folders",
    "attributes": {
      "name": "Test Folder",
      "user_id": "{user_id}"
    }
  }
}
```

#### Create a Subfolder

- **URL**: `/v1/folders`
- **Method**: `POST`
- **Request Body**:
```json
{
  "data": {
    "type": "folders",
    "attributes": {
      "name": "Test Subfolder",
      "user_id": "{user_id}",
      "parent_id": "{parent_folder_id}"
    }
  }
}
```

#### Get All Folders

- **URL**: `/v1/folders`
- **Method**: `GET`

#### Get Folders by User ID

- **URL**: `/v1/folders?user_id={user_id}`
- **Method**: `GET`

#### Get Subfolders by Parent ID

- **URL**: `/v1/folders?parent_id={parent_folder_id}`
- **Method**: `GET`

#### Get Root Folders (no parent)

- **URL**: `/v1/folders?parent_id=null`
- **Method**: `GET`

#### Get a Folder

- **URL**: `/v1/folders/{id}`
- **Method**: `GET`

#### Update a Folder

- **URL**: `/v1/folders/{id}`
- **Method**: `PATCH`
- **Request Body**:
```json
{
  "data": {
    "type": "folders",
    "id": "{id}",
    "attributes": {
      "name": "Updated Folder Name"
    }
  }
}
```

#### Move a Folder to Another Parent

- **URL**: `/v1/folders/{id}`
- **Method**: `PATCH`
- **Request Body**:
```json
{
  "data": {
    "type": "folders",
    "id": "{id}",
    "attributes": {
      "parent_id": "{new_parent_folder_id}"
    }
  }
}
```

#### Delete a Folder

- **URL**: `/v1/folders/{id}`
- **Method**: `DELETE`

### Documents

#### Create a Document

- **URL**: `/v1/documents`
- **Method**: `POST`
- **Request Body**:
```json
{
  "data": {
    "type": "documents",
    "attributes": {
      "title": "Test Document",
      "content": "This is a test document content.",
      "user_id": "{user_id}"
    }
  }
}
```

#### Create a Document in a Folder

- **URL**: `/v1/documents`
- **Method**: `POST`
- **Request Body**:
```json
{
  "data": {
    "type": "documents",
    "attributes": {
      "title": "Test Document in Folder",
      "content": "This is a test document content in a folder.",
      "user_id": "{user_id}",
      "folder_id": "{folder_id}"
    }
  }
}
```

#### Get All Documents

- **URL**: `/v1/documents`
- **Method**: `GET`

#### Get Documents by User ID

- **URL**: `/v1/documents?user_id={user_id}`
- **Method**: `GET`

#### Get Documents by Folder ID

- **URL**: `/v1/documents?folder_id={folder_id}`
- **Method**: `GET`

#### Get Documents with No Folder

- **URL**: `/v1/documents?folder_id=null`
- **Method**: `GET`

#### Get a Document

- **URL**: `/v1/documents/{id}`
- **Method**: `GET`

#### Update a Document

- **URL**: `/v1/documents/{id}`
- **Method**: `PATCH`
- **Request Body**:
```json
{
  "data": {
    "type": "documents",
    "id": "{id}",
    "attributes": {
      "title": "Updated Document Title",
      "content": "This is the updated document content."
    }
  }
}
```

#### Move a Document to a Folder

- **URL**: `/v1/documents/{id}`
- **Method**: `PATCH`
- **Request Body**:
```json
{
  "data": {
    "type": "documents",
    "id": "{id}",
    "attributes": {
      "folder_id": "{folder_id}"
    }
  }
}
```

#### Remove a Document from a Folder

- **URL**: `/v1/documents/{id}`
- **Method**: `PATCH`
- **Request Body**:
```json
{
  "data": {
    "type": "documents",
    "id": "{id}",
    "attributes": {
      "folder_id": null
    }
  }
}
```

#### Delete a Document

- **URL**: `/v1/documents/{id}`
- **Method**: `DELETE`

## Testing with Bruno

The project includes Bruno API definitions for testing the endpoints. To use them:

1. Install Bruno: https://www.usebruno.com/
2. Open Bruno and import the `.bruno` directory
3. Set up the environment variables in `.bruno/environments/local.bru`
4. Run the requests to test the API

## License

This project is licensed under the MIT License - see the LICENSE file for details.
