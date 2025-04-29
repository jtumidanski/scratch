// API Types based on the JSON:API specification

// Base Types
export type ApiResource = {
  id: string;
  type: string;
  attributes: Record<string, any>;
  relationships?: Record<string, ApiRelationship>;
};

export type ApiRelationship = {
  data: { id: string; type: string } | { id: string; type: string }[] | null;
};

export type ApiResponse<T extends ApiResource> = {
  data: T | T[];
  included?: ApiResource[];
  meta?: Record<string, any>;
};

// User Types
export type UserAttributes = {
  username: string;
  email: string;
};

export type User = ApiResource & {
  type: 'users';
  attributes: UserAttributes;
};

// Folder Types
export type FolderAttributes = {
  name: string;
  user_id: string;
  parent_id?: string | null;
};

export type Folder = ApiResource & {
  type: 'folders';
  attributes: FolderAttributes;
};

// Document Types
export type DocumentAttributes = {
  title: string;
  content: string;
  user_id: string;
  folder_id?: string | null;
};

export type Document = ApiResource & {
  type: 'documents';
  attributes: DocumentAttributes;
};

// Request Types
export type CreateUserRequest = {
  data: {
    type: 'users';
    id: string;
    attributes: {
      username: string;
      email: string;
    };
  };
};

export type UpdateUserRequest = {
  data: {
    type: 'users';
    id: string;
    attributes: Partial<UserAttributes>;
  };
};

export type CreateFolderRequest = {
  data: {
    type: 'folders';
    id: string;
    attributes: {
      name: string;
      user_id: string;
      parent_id?: string | null;
    };
  };
};

export type UpdateFolderRequest = {
  data: {
    type: 'folders';
    id: string;
    attributes: Partial<FolderAttributes>;
  };
};

export type CreateDocumentRequest = {
  data: {
    type: 'documents';
    id: string;
    attributes: {
      title: string;
      content: string;
      user_id: string;
      folder_id?: string | null;
    };
  };
};

export type UpdateDocumentRequest = {
  data: {
    type: 'documents';
    id: string;
    attributes: Partial<DocumentAttributes>;
  };
};
