import {
  ApiResponse,
  User,
  Folder,
  Document,
  CreateUserRequest,
  UpdateUserRequest,
  CreateFolderRequest,
  UpdateFolderRequest,
  CreateDocumentRequest,
  UpdateDocumentRequest
} from './api-types';
import { NIL_UUID } from './utils';

// Base API URL from environment variables with fallback
// Using typeof window to check if we're in browser context
// This ensures we get the latest value of the environment variable at runtime
const API_BASE_URL = process.env.NEXT_PUBLIC_ROOT_API_URL || window.location.origin;

// Generic fetch function with error handling
async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {})
  };

  const response = await fetch(url, {
    ...options,
    headers
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.message || `API error: ${response.status}`);
  }

  // If the response is 204 No Content, return undefined instead of trying to parse JSON
  if (response.status === 204) {
    return undefined as unknown as T;
  }

  return response.json();
}

// User API functions
export async function getUsers(): Promise<User[]> {
  const response = await fetchApi<ApiResponse<User>>('/v1/users');
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getUser(id: string): Promise<User> {
  const response = await fetchApi<ApiResponse<User>>(`/v1/users/${id}`);
  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function getUserByUsername(username: string): Promise<User | null> {
  try {
    const users = await getUsers();
    const user = users.find(u => u.attributes.username === username);
    return user || null;
  } catch (error) {
    console.error('Error fetching user by username:', error);
    return null;
  }
}

export async function getUserOrCreate(username: string, email?: string): Promise<User | null> {
  try {
    // Try to find the user by username
    const user = await getUserByUsername(username);

    // If user exists, return it
    if (user) {
      return user;
    }

    // If user doesn't exist and we have an email, create a new user
    if (email) {
      return await createUser(username, email);
    }

    // If no email is provided, we can't create a user
    console.error('Cannot create user: email is required');
    return null;
  } catch (error) {
    console.error('Error in getUserOrCreate:', error);
    return null;
  }
}

export async function createUser(username: string, email: string): Promise<User> {
  const request: CreateUserRequest = {
    data: {
      type: 'users',
      id: NIL_UUID,
      attributes: {
        username,
        email
      }
    }
  };

  const response = await fetchApi<ApiResponse<User>>('/v1/users', {
    method: 'POST',
    body: JSON.stringify(request)
  });

  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function updateUser(id: string, data: Partial<{ username: string, email: string }>): Promise<User> {
  const request: UpdateUserRequest = {
    data: {
      type: 'users',
      id,
      attributes: data
    }
  };

  const response = await fetchApi<ApiResponse<User>>(`/v1/users/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(request)
  });

  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function deleteUser(id: string): Promise<void> {
  await fetchApi(`/users/${id}`, {
    method: 'DELETE'
  });
}

// Folder API functions
export async function getFolders(): Promise<Folder[]> {
  const response = await fetchApi<ApiResponse<Folder>>('/v1/folders');
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getFoldersByUserId(userId: string): Promise<Folder[]> {
  const response = await fetchApi<ApiResponse<Folder>>(`/v1/folders?user_id=${userId}`);
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getSubfoldersByParentId(parentId: string): Promise<Folder[]> {
  const response = await fetchApi<ApiResponse<Folder>>(`/v1/folders?parent_id=${parentId}`);
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getRootFolders(): Promise<Folder[]> {
  const response = await fetchApi<ApiResponse<Folder>>('/v1/folders?parent_id=null');
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getFolder(id: string): Promise<Folder> {
  const response = await fetchApi<ApiResponse<Folder>>(`/v1/folders/${id}`);
  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function createFolder(name: string, userId: string, parentId?: string): Promise<Folder> {
  const request: CreateFolderRequest = {
    data: {
      type: 'folders',
      id: NIL_UUID,
      attributes: {
        name,
        user_id: userId,
        ...(parentId && { parent_id: parentId })
      }
    }
  };

  const response = await fetchApi<ApiResponse<Folder>>('/v1/folders', {
    method: 'POST',
    body: JSON.stringify(request)
  });

  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function updateFolder(id: string, data: Partial<{ name: string, parent_id?: string | null }>): Promise<Folder> {
  const request: UpdateFolderRequest = {
    data: {
      type: 'folders',
      id,
      attributes: data
    }
  };

  const response = await fetchApi<ApiResponse<Folder>>(`/v1/folders/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(request)
  });

  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function deleteFolder(id: string): Promise<void> {
  await fetchApi(`/folders/${id}`, {
    method: 'DELETE'
  });
}

// Document API functions
export async function getDocuments(): Promise<Document[]> {
  const response = await fetchApi<ApiResponse<Document>>('/v1/documents');
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getDocumentsByUserId(userId: string): Promise<Document[]> {
  const response = await fetchApi<ApiResponse<Document>>(`/v1/documents?user_id=${userId}`);
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getDocumentsByFolderId(folderId: string): Promise<Document[]> {
  const response = await fetchApi<ApiResponse<Document>>(`/v1/documents?folder_id=${folderId}`);
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getDocumentsWithNoFolder(): Promise<Document[]> {
  const response = await fetchApi<ApiResponse<Document>>('/v1/documents?folder_id=null');
  return Array.isArray(response.data) ? response.data : [response.data];
}

export async function getDocument(id: string): Promise<Document> {
  const response = await fetchApi<ApiResponse<Document>>(`/v1/documents/${id}`);
  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function createDocument(title: string, content: string, userId: string, folderId?: string): Promise<Document> {
  const request: CreateDocumentRequest = {
    data: {
      type: 'documents',
      id: NIL_UUID,
      attributes: {
        title,
        content,
        user_id: userId,
        ...(folderId && { folder_id: folderId })
      }
    }
  };

  const response = await fetchApi<ApiResponse<Document>>('/v1/documents', {
    method: 'POST',
    body: JSON.stringify(request)
  });

  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function updateDocument(id: string, data: Partial<{ title: string, content: string, folder_id?: string | null }>): Promise<Document> {
  const request: UpdateDocumentRequest = {
    data: {
      type: 'documents',
      id,
      attributes: data
    }
  };

  const response = await fetchApi<ApiResponse<Document>>(`/v1/documents/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(request)
  });

  return Array.isArray(response.data) ? response.data[0] : response.data;
}

export async function deleteDocument(id: string): Promise<void> {
  await fetchApi(`/v1/documents/${id}`, {
    method: 'DELETE'
  });
}
