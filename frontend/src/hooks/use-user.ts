import { useQuery } from "@tanstack/react-query"

export type UserProfile = {
  id: number
  name: string
  username: string
  email: string
  image?: string
}

export function useUser() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["user-profile"],
    queryFn: async () => {
      const res = await fetch("/api/user/profile")
      if (!res.ok) {
        throw new Error("Failed to fetch user profile")
      }
      const body = await res.json()
      return body.data as UserProfile
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 1,
  })

  return {
    user: data,
    isLoading,
    error,
    refetch,
  }
}
