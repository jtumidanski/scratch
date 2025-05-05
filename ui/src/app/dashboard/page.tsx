'use client'

import { useState, useEffect, useRef, useCallback } from 'react'
import { DeleteDocumentDialog } from '@/components/documents/delete-document-dialog'
import { DeleteFolderDialog } from '@/components/folders/delete-folder-dialog'
import { CreateItemDialog } from '@/components/create-item-dialog'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Folder as FolderIcon, FolderOpen, FileText, Trash2, Plus } from 'lucide-react'
import { 
  getUser, 
  getUsers,
  getUserByUsername,
  getUserOrCreate,
  getFoldersByUserId, 
  getDocumentsByUserId, 
  updateDocument,
  deleteDocument,
  deleteFolder
} from '@/lib/api'
import { User, Folder, Document } from '@/lib/api-types'

// Custom debounce function implementation
const debounce = (func: Function, wait: number) => {
  let timeout: NodeJS.Timeout | null = null;

  return (...args: any[]) => {
    if (timeout) {
      clearTimeout(timeout);
    }

    timeout = setTimeout(() => {
      func(...args);
    }, wait);
  };
};

// Define types for navigation
type ItemType = 'folder' | 'document';
type NavigationItem = {
  id: string;
  type: ItemType;
  parentId: string | null;
  element: HTMLElement | null;
};

// Define type for history entries
type HistoryEntry = {
  content: string;
  cursorPosition: {
    start: number;
    end: number;
  };
};

// Configuration for edit history
const HISTORY_SIZE = 50;

export default function DashboardPage() {
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [folders, setFolders] = useState<Folder[]>([])
  const [documents, setDocuments] = useState<Document[]>([])
  const [selectedDocument, setSelectedDocument] = useState<Document | null>(null)
  const [localContent, setLocalContent] = useState<string>('')
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set())
  const [isLoading, setIsLoading] = useState(true)
  const [currentFolderId, setCurrentFolderId] = useState<string | undefined>(undefined)

  // Edit history state
  const [editHistory, setEditHistory] = useState<HistoryEntry[]>([])
  const [historyPosition, setHistoryPosition] = useState(-1)
  const [isUndoRedoAction, setIsUndoRedoAction] = useState(false)

  // Navigation state
  const [focusedItemId, setFocusedItemId] = useState<string | null>(null)
  const [navigationItems, setNavigationItems] = useState<NavigationItem[]>([])
  const itemRefs = useRef<Map<string, HTMLElement>>(new Map())
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // Dialog states
  const [deleteDocumentDialogOpen, setDeleteDocumentDialogOpen] = useState(false)
  const [documentToDelete, setDocumentToDelete] = useState<{ id: string, title: string } | null>(null)
  const [deleteFolderDialogOpen, setDeleteFolderDialogOpen] = useState(false)
  const [folderToDelete, setFolderToDelete] = useState<{ id: string, name: string } | null>(null)

  const fetchData = async () => {
    setIsLoading(true)
    try {
      // Get the username from localStorage
      const username = localStorage.getItem('username')
      if (!username) return

      // Try to get the user or create if not found
      let user = await getUserByUsername(username)

      // If user not found, prompt for email and create user
      if (!user) {
        const email = prompt(`User "${username}" not found. Please enter your email to create a new account:`)

        // If email provided, create the user
        if (email) {
          user = await getUserOrCreate(username, email)
        } else {
          console.error('Email required to create user')
          return
        }

        // If still no user (creation failed), return
        if (!user) {
          console.error('Failed to create user')
          return
        }
      }

      setCurrentUser(user)

      // Fetch folders and documents from the API using the user's ID
      const userFolders = await getFoldersByUserId(user.id)
      const userDocuments = await getDocumentsByUserId(user.id)

      setFolders(userFolders)
      setDocuments(userDocuments)
    } catch (error) {
      console.error('Error fetching data:', error)
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // Build the navigation items list whenever folders, documents, or expandedFolders change
  useEffect(() => {
    if (isLoading) return;

    const items: NavigationItem[] = [];

    // Add root documents
    documents
      .filter(doc => doc.attributes.folder_id === null)
      .forEach(doc => {
        items.push({
          id: doc.id,
          type: 'document',
          parentId: null,
          element: itemRefs.current.get(doc.id) || null
        });
      });

    // Function to recursively add folders and their contents
    const addFolderItems = (parentId: string | null) => {
      const folderItems = folders.filter(folder => folder.attributes.parent_id === parentId);

      folderItems.forEach(folder => {
        // Add the folder
        items.push({
          id: folder.id,
          type: 'folder',
          parentId,
          element: itemRefs.current.get(folder.id) || null
        });

        // If folder is expanded, add its contents
        if (expandedFolders.has(folder.id)) {
          // Recursively add subfolders
          addFolderItems(folder.id);

          // Add documents in this folder
          documents
            .filter(doc => doc.attributes.folder_id === folder.id)
            .forEach(doc => {
              items.push({
                id: doc.id,
                type: 'document',
                parentId: folder.id,
                element: itemRefs.current.get(doc.id) || null
              });
            });
        }
      });
    };

    // Start with root folders
    addFolderItems(null);

    setNavigationItems(items);
  }, [folders, documents, expandedFolders, isLoading])

  const toggleFolder = (folderId: string) => {
    setExpandedFolders(prev => {
      const newSet = new Set(prev)
      if (newSet.has(folderId)) {
        newSet.delete(folderId)
      } else {
        newSet.add(folderId)
      }
      return newSet
    })
  }

  const handleDocumentSelect = (document: Document) => {
    setSelectedDocument(document)
    const content = document.attributes.content
    setLocalContent(content)

    // Reset edit history when selecting a new document
    // Initialize with default cursor position at the start of the document
    setEditHistory([{
      content,
      cursorPosition: {
        start: 0,
        end: 0
      }
    }])
    setHistoryPosition(0)

    // Use setTimeout to ensure the Textarea is rendered before trying to focus it
    setTimeout(() => {
      if (textareaRef.current) {
        textareaRef.current.focus()
      }
    }, 0)
  }

  // Global keyboard navigation
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      // Skip if a dialog is open
      if (deleteDocumentDialogOpen || deleteFolderDialogOpen) return;

      // Skip if the active element is the Textarea (main content box)
      const activeElement = document.activeElement;
      if (activeElement && activeElement.tagName === 'TEXTAREA') return;

      // Skip if the active element is an input field
      if (activeElement && (
        activeElement.tagName === 'INPUT' || 
        activeElement.tagName === 'SELECT' || 
        activeElement.tagName === 'BUTTON'
      )) return;

      const currentIndex = focusedItemId 
        ? navigationItems.findIndex(item => item.id === focusedItemId)
        : -1;

      // Handle Enter and Space keys for selection
      if ((e.key === 'Enter' || e.key === ' ') && currentIndex !== -1) {
        e.preventDefault();
        const currentItem = navigationItems[currentIndex];
        if (currentItem.type === 'folder') {
          toggleFolder(currentItem.id);
        } else if (currentItem.type === 'document') {
          const doc = documents.find(d => d.id === currentItem.id);
          if (doc) {
            handleDocumentSelect(doc);
          }
        }
        return;
      }

      // Handle arrow keys for navigation
      if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(e.key)) {
        e.preventDefault();

        switch (e.key) {
          case 'ArrowUp':
            if (currentIndex > 0) {
              const prevItem = navigationItems[currentIndex - 1];
              setFocusedItemId(prevItem.id);
              prevItem.element?.focus();
            }
            break;

          case 'ArrowDown':
            if (currentIndex < navigationItems.length - 1) {
              const nextItem = navigationItems[currentIndex + 1];
              setFocusedItemId(nextItem.id);
              nextItem.element?.focus();
            } else if (currentIndex === -1 && navigationItems.length > 0) {
              // If nothing is focused yet, focus the first item
              const firstItem = navigationItems[0];
              setFocusedItemId(firstItem.id);
              firstItem.element?.focus();
            }
            break;

          case 'ArrowLeft':
            if (currentIndex !== -1) {
              const currentItem = navigationItems[currentIndex];
              if (currentItem.type === 'folder' && expandedFolders.has(currentItem.id)) {
                // Collapse the folder if it's expanded
                toggleFolder(currentItem.id);
              } else if (currentItem.parentId) {
                // Navigate to parent folder
                const parentIndex = navigationItems.findIndex(item => item.id === currentItem.parentId);
                if (parentIndex !== -1) {
                  const parentItem = navigationItems[parentIndex];
                  setFocusedItemId(parentItem.id);
                  parentItem.element?.focus();
                }
              }
            }
            break;

          case 'ArrowRight':
            if (currentIndex !== -1) {
              const currentItem = navigationItems[currentIndex];
              if (currentItem.type === 'folder' && !expandedFolders.has(currentItem.id)) {
                // Expand the folder if it's collapsed
                toggleFolder(currentItem.id);
              } else if (currentItem.type === 'folder' && expandedFolders.has(currentItem.id)) {
                // If folder is already expanded, move to its first child
                const nextIndex = currentIndex + 1;
                if (nextIndex < navigationItems.length) {
                  const nextItem = navigationItems[nextIndex];
                  if (nextItem.parentId === currentItem.id) {
                    setFocusedItemId(nextItem.id);
                    nextItem.element?.focus();
                  }
                }
              }
            }
            break;
        }
      }
    };

    // Add event listener
    window.addEventListener('keydown', handleGlobalKeyDown);

    // Remove event listener on cleanup
    return () => {
      window.removeEventListener('keydown', handleGlobalKeyDown);
    };
  }, [navigationItems, focusedItemId, expandedFolders, toggleFolder, setFocusedItemId, deleteDocumentDialogOpen, deleteFolderDialogOpen, documents, handleDocumentSelect])

  // Create a debounced function for API updates
  const debouncedUpdateDocument = useCallback(
    debounce(async (documentId: string, content: string) => {
      try {
        const updatedDocument = await updateDocument(documentId, { content });

        setSelectedDocument(updatedDocument);

        // Update the documents array with the updated document
        setDocuments(prev => 
          prev.map(doc => 
            doc.id === documentId 
              ? updatedDocument
              : doc
          )
        );
      } catch (error) {
        console.error('Error updating document:', error);
      }
    }, 500), // 500ms debounce time
    []
  );

  // Undo function - go back one step in history
  const handleUndo = () => {
    if (historyPosition > 0 && textareaRef.current) {
      setIsUndoRedoAction(true);
      const newPosition = historyPosition - 1;
      const historyEntry = editHistory[newPosition];

      setHistoryPosition(newPosition);
      setLocalContent(historyEntry.content);

      // Restore the cursor position that was saved with this history entry
      setTimeout(() => {
        if (textareaRef.current) {
          textareaRef.current.selectionStart = historyEntry.cursorPosition.start;
          textareaRef.current.selectionEnd = historyEntry.cursorPosition.end;
          textareaRef.current.focus();
        }
      }, 0);

      // Debounce the API update with the content from history
      if (selectedDocument) {
        debouncedUpdateDocument(selectedDocument.id, historyEntry.content);
      }
    }
  };

  // Redo function - go forward one step in history
  const handleRedo = () => {
    if (historyPosition < editHistory.length - 1 && textareaRef.current) {
      setIsUndoRedoAction(true);
      const newPosition = historyPosition + 1;
      const historyEntry = editHistory[newPosition];

      setHistoryPosition(newPosition);
      setLocalContent(historyEntry.content);

      // Restore the cursor position that was saved with this history entry
      setTimeout(() => {
        if (textareaRef.current) {
          textareaRef.current.selectionStart = historyEntry.cursorPosition.start;
          textareaRef.current.selectionEnd = historyEntry.cursorPosition.end;
          textareaRef.current.focus();
        }
      }, 0);

      // Debounce the API update with the content from history
      if (selectedDocument) {
        debouncedUpdateDocument(selectedDocument.id, historyEntry.content);
      }
    }
  };


  const handleDocumentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    if (selectedDocument) {
      const newContent = e.target.value;
      const textarea = e.target;

      // Update local state immediately (this preserves cursor position)
      setLocalContent(newContent);

      // If this change is from an undo/redo action, don't add to history
      if (isUndoRedoAction) {
        setIsUndoRedoAction(false);
        return;
      }

      // Create a new history entry with content and cursor position
      const newEntry: HistoryEntry = {
        content: newContent,
        cursorPosition: {
          start: textarea.selectionStart,
          end: textarea.selectionEnd
        }
      };

      // Add to history, truncating any future history if we're not at the end
      // This ensures that undone edits are removed from the cache when new edits are made
      const newHistory = [...editHistory.slice(0, historyPosition + 1), newEntry];

      // Limit history size
      if (newHistory.length > HISTORY_SIZE) {
        newHistory.shift(); // Remove oldest entry
      }

      setEditHistory(newHistory);
      setHistoryPosition(newHistory.length - 1);

      // Debounce the API update
      debouncedUpdateDocument(selectedDocument.id, newContent);
    }
  }

  const handleDocumentDelete = (document: Document) => {
    setDocumentToDelete({ id: document.id, title: document.attributes.title })
    setDeleteDocumentDialogOpen(true)
  }

  const handleDocumentDeleted = (documentId: string) => {
    // Update the local state
    setDocuments(prev => prev.filter(doc => doc.id !== documentId))

    if (selectedDocument?.id === documentId) {
      setSelectedDocument(null)
    }
  }

  const handleFolderDelete = (folder: Folder) => {
    setFolderToDelete({ id: folder.id, name: folder.attributes.name })
    setDeleteFolderDialogOpen(true)
  }

  const handleFolderDeleted = (folderId: string) => {
    // Update the local state
    setFolders(prev => prev.filter(folder => folder.id !== folderId))

    // Move documents to root in the local state
    // The API should handle this on the server side, but we'll update our local state as well
    setDocuments(prev => 
      prev.map(doc => 
        doc.attributes.folder_id === folderId 
          ? { ...doc, attributes: { ...doc.attributes, folder_id: null } }
          : doc
      )
    )
  }

  const handleFolderSelect = (folderId: string | undefined) => {
    setCurrentFolderId(folderId)
  }


  // Recursive function to render folders and their subfolders
  const renderFolderTree = (parentId: string | null) => {
    const folderItems = folders.filter(folder => folder.attributes.parent_id === parentId)

    return folderItems.map(folder => {
      const isExpanded = expandedFolders.has(folder.id)
      const folderDocuments = documents.filter(doc => doc.attributes.folder_id === folder.id)

      return (
        <div key={folder.id}>
          <div className={`flex items-center justify-between group hover:bg-muted ${focusedItemId === folder.id ? 'bg-primary/20 outline-none ring-2 ring-primary' : ''}`}>
            <div 
              className="flex items-center p-2 cursor-pointer flex-grow select-none"
              onClick={() => toggleFolder(folder.id)}
              ref={(el) => {
                if (el) {
                  itemRefs.current.set(folder.id, el);
                }
              }}
            >
              {isExpanded ? <FolderOpen className="mr-1 h-4 w-4" /> : <FolderIcon className="mr-1 h-4 w-4" />}
              <span>{folder.attributes.name}</span>
            </div>
            <div className="opacity-0 group-hover:opacity-100 transition-opacity pr-2 flex">
              {currentUser && (
                <div className="h-6 w-6 flex items-center justify-center">
                  <CreateItemDialog 
                    userId={currentUser.id} 
                    parentId={folder.id}
                    onItemCreated={fetchData}
                    triggerLabel=""
                  />
                </div>
              )}
              <Button 
                variant="ghost" 
                size="sm" 
                className="h-6 w-6 p-0 ml-1 group-hover:text-foreground" 
                onClick={(e) => {
                  e.stopPropagation()
                  handleFolderDelete(folder)
                }}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {isExpanded && (
            <div className="ml-4">
              {renderFolderTree(folder.id)}

              {folderDocuments.map(doc => (
                <div 
                  key={doc.id}
                  className={`flex items-center justify-between group ${selectedDocument?.id === doc.id ? 'bg-primary/10' : 'hover:bg-muted'} ${focusedItemId === doc.id ? 'bg-primary/20 outline-none ring-2 ring-primary' : ''}`}
                >
                  <div 
                    className="flex items-center p-2 cursor-pointer flex-grow select-none"
                    onClick={() => handleDocumentSelect(doc)}
                    ref={(el) => {
                      if (el) {
                        itemRefs.current.set(doc.id, el);
                      }
                    }}
                  >
                    <FileText className="mr-1 h-4 w-4" />
                    <span>{doc.attributes.title}</span>
                  </div>
                  <div className="opacity-0 group-hover:opacity-100 transition-opacity pr-2">
                    <Button 
                      variant="ghost" 
                      size="sm" 
                      className="h-6 w-6 p-0" 
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDocumentDelete(doc)
                      }}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )
    })
  }

  if (isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <p className="text-muted-foreground">Loading documents...</p>
      </div>
    )
  }

  return (
    <>
      {/* Left segment - Document/Folder Navigation */}
      <div 
        className="w-1/4 border-r overflow-y-auto p-4"
        onKeyDown={(e) => {
          // Only handle arrow down to start navigation if nothing is focused
          if (e.key === 'ArrowDown' && !focusedItemId && navigationItems.length > 0) {
            e.preventDefault();
            const firstItem = navigationItems[0];
            setFocusedItemId(firstItem.id);
            firstItem.element?.focus();
          }
        }}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold">Documents</h2>
          <div className="flex space-x-2">
            {currentUser && (
              <CreateItemDialog 
                userId={currentUser.id} 
                parentId={currentFolderId}
                onItemCreated={fetchData} 
              />
            )}
          </div>
        </div>

        {/* Root level documents (no folder) */}
        {documents
          .filter(doc => doc.attributes.folder_id === null)
          .map(doc => (
            <div 
              key={doc.id}
              className={`flex items-center justify-between group ${selectedDocument?.id === doc.id ? 'bg-primary/10' : 'hover:bg-muted'} ${focusedItemId === doc.id ? 'bg-primary/20 outline-none ring-2 ring-primary' : ''}`}
            >
              <div 
                className="flex items-center p-2 cursor-pointer flex-grow select-none"
                onClick={() => handleDocumentSelect(doc)}
                ref={(el) => {
                  if (el) {
                    itemRefs.current.set(doc.id, el);
                  }
                }}
              >
                <FileText className="mr-1 h-4 w-4" />
                <span>{doc.attributes.title}</span>
              </div>
              <div className="opacity-0 group-hover:opacity-100 transition-opacity pr-2">
                <Button 
                  variant="ghost" 
                  size="sm" 
                  className="h-6 w-6 p-0" 
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDocumentDelete(doc)
                  }}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))
        }

        {/* Folder tree */}
        {renderFolderTree(null)}
      </div>

      {/* Right segment - Document Content */}
      <div className="w-3/4 p-4 overflow-y-auto">
        {selectedDocument ? (
          <div className="h-full flex flex-col">
            <h2 className="text-xl font-semibold mb-4">{selectedDocument.attributes.title}</h2>
            <Textarea
              ref={textareaRef}
              className="flex-1 min-h-[300px]"
              value={localContent}
              onChange={handleDocumentChange}
              onKeyDown={(e) => {
                if (e.key === 'Escape') {
                  e.currentTarget.blur();
                }

                // Undo: Ctrl+Z
                if ((e.ctrlKey || e.metaKey) && e.key === 'z' && !e.shiftKey) {
                  e.preventDefault();
                  handleUndo();
                }

                // Redo: Ctrl+Y or Ctrl+Shift+Z
                if ((e.ctrlKey || e.metaKey) && (e.key === 'y' || (e.key === 'z' && e.shiftKey))) {
                  e.preventDefault();
                  handleRedo();
                }
              }}
            />
          </div>
        ) : (
          <div className="h-full flex items-center justify-center text-muted-foreground">
            <p>Select a document to view and edit</p>
          </div>
        )}
      </div>

      {/* Delete Confirmation Dialogs */}
      {documentToDelete && (
        <DeleteDocumentDialog
          documentId={documentToDelete.id}
          documentTitle={documentToDelete.title}
          onDocumentDeleted={handleDocumentDeleted}
          open={deleteDocumentDialogOpen}
          onOpenChange={setDeleteDocumentDialogOpen}
        />
      )}

      {folderToDelete && (
        <DeleteFolderDialog
          folderId={folderToDelete.id}
          folderName={folderToDelete.name}
          onFolderDeleted={handleFolderDeleted}
          open={deleteFolderDialogOpen}
          onOpenChange={setDeleteFolderDialogOpen}
        />
      )}
    </>
  )
}
