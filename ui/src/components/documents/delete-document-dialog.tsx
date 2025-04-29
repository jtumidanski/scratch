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
import { deleteDocument } from "@/lib/api"

interface DeleteDocumentDialogProps {
  documentId: string
  documentTitle: string
  onDocumentDeleted: (documentId: string) => void
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DeleteDocumentDialog({ 
  documentId, 
  documentTitle,
  onDocumentDeleted,
  open,
  onOpenChange
}: DeleteDocumentDialogProps) {
  const [isLoading, setIsLoading] = useState(false)

  const handleDelete = async () => {
    setIsLoading(true)
    try {
      // Delete the document via the API
      await deleteDocument(documentId)
      onDocumentDeleted(documentId)
      onOpenChange(false)
    } catch (error) {
      console.error('Error deleting document:', error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Document</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete "{documentTitle}"? This action cannot be undone.
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