"use client"

import { useState } from "react"
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter,
  DialogTrigger
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { createDocument, createFolder } from "@/lib/api"
import { Plus } from "lucide-react"

interface CreateItemDialogProps {
  userId: string
  parentId?: string
  onItemCreated: () => void
  triggerLabel?: string
}

type ItemType = "document" | "folder"

export function CreateItemDialog({ 
  userId, 
  parentId, 
  onItemCreated,
  triggerLabel = "New Item"
}: CreateItemDialogProps) {
  const [name, setName] = useState("")
  const [itemType, setItemType] = useState<ItemType>("document")
  const [isLoading, setIsLoading] = useState(false)
  const [isOpen, setIsOpen] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) return

    setIsLoading(true)
    try {
      if (itemType === "document") {
        await createDocument(name, "", userId, parentId)
      } else {
        await createFolder(name, userId, parentId)
      }
      setName("")
      setIsOpen(false)
      onItemCreated()
    } catch (error) {
      console.error(`Failed to create ${itemType}:`, error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button 
          variant={triggerLabel ? "outline" : "ghost"} 
          size="sm" 
          className={`flex items-center ${!triggerLabel ? 'h-6 w-6 p-0 group-hover:text-foreground' : ''}`}
        >
          <Plus className={`h-4 w-4 ${triggerLabel ? 'mr-1' : ''}`} />
          {triggerLabel}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create New {itemType === "document" ? "Document" : "Folder"}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="itemType" className="text-right">
                Type
              </Label>
              <div className="col-span-3 flex space-x-2">
                <Button 
                  type="button"
                  variant={itemType === "document" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setItemType("document")}
                >
                  Document
                </Button>
                <Button 
                  type="button"
                  variant={itemType === "folder" ? "default" : "outline"}
                  size="sm"
                  onClick={() => setItemType("folder")}
                >
                  Folder
                </Button>
              </div>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right">
                {itemType === "document" ? "Title" : "Name"}
              </Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="col-span-3"
                placeholder={`Enter ${itemType} ${itemType === "document" ? "title" : "name"}`}
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Creating..." : "Create"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
