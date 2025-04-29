'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  useEffect(() => {
    // Check if user is authenticated
    const username = localStorage.getItem('username')
    if (!username) {
      router.push('/')
    } else {
      setIsLoading(false)
    }
  }, [router])

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <p className="text-lg text-muted-foreground">Loading...</p>
      </div>
    )
  }

  return (
    <div className="flex h-screen flex-col">
      <header className="bg-primary text-primary-foreground p-4 border-b border-border shadow-sm">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-xl font-bold">Scratch Document Service</h1>
          <button
            onClick={() => {
              localStorage.removeItem('username')
              router.push('/')
            }}
            className="px-3 py-1 bg-accent hover:bg-accent/90 text-accent-foreground rounded transition-colors"
          >
            Sign Out
          </button>
        </div>
      </header>
      <main className="flex flex-1 overflow-hidden">
        {children}
      </main>
    </div>
  )
}
