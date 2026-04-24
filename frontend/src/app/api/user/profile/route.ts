import { cookies } from "next/headers"
import { NextResponse } from "next/server"

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"

export async function GET() {
  const cookieStore = await cookies()
  const token = cookieStore.get("auth_token")?.value

  if (!token) {
    return NextResponse.json(
      { code: "ERR_UNAUTHORIZED", message: "Missing authentication token" },
      { status: 401 }
    )
  }

  try {
    const res = await fetch(`${API_URL}/api/v1/profile`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      cache: "no-store",
    })

    const data = await res.json()

    if (!res.ok) {
      return NextResponse.json(data, { status: res.status })
    }

    return NextResponse.json(data)
  } catch (error) {
    console.error("Proxy Profile Error:", error)
    return NextResponse.json(
      { code: "ERR_INTERNAL", message: "Failed to connect to backend" },
      { status: 500 }
    )
  }
}
