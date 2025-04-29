"use client"

import { useState } from "react"
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter,
  DialogDescription
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { deleteFolder } from "@/lib/api"

interface DeleteFolderDialogProps {
  folderId: string
  folderName: string
  onFolderDeleted: (folderId: string) => void
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DeleteFolderDialog({ 
  folderId, 
  folderName,
  onFolderDeleted,
  open,
  onOpenChange
}: DeleteFolderDialogProps) {
  const [isLoading, setIsLoading] = useState(false)

  const handleDelete = async () => {
    setIsLoading(true)
    try {
      // Delete the folder via the API
      await deleteFolder(folderId)
      onFolderDeleted(folderId)
      onOpenChange(false)
    } catch (error) {
      console.error('Error deleting folder:', error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Folder</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete "{folderName}"? All documents inside will be moved to root. This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={isLoading}>
            {isLoading ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}