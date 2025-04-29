import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

// Nil UUID - a UUID with all zeros
export const NIL_UUID = "00000000-0000-0000-0000-000000000000"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
