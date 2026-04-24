import { cookies } from "next/headers"
import { NextRequest, NextResponse } from "next/server"

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

async function proxyRequest(request: NextRequest) {
  const cookieStore = await cookies()
  const token = cookieStore.get("auth_token")?.value

  // Extract the path from the URL
  const { pathname, search } = new URL(request.url)
  
  // Construct the backend URL
  // The path will already start with /api/v1 because of the folder structure
  const backendUrl = `${API_URL}${pathname}${search}`

  const headers = new Headers(request.headers)
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  // Next.js might block some headers or add its own, we want to be clean
  headers.delete("host")
  headers.delete("connection")

  try {
    const body = request.method !== 'GET' && request.method !== 'HEAD' 
      ? await request.arrayBuffer() 
      : undefined

    const res = await fetch(backendUrl, {
      method: request.method,
      headers: headers,
      body: body,
      cache: "no-store",
    })

    const contentType = res.headers.get("content-type")
    let data
    
    if (contentType?.includes("application/json")) {
      data = await res.json()
    } else {
      data = await res.text()
    }

    // Return the response directly
    return NextResponse.json(data, {
      status: res.status,
      headers: {
        "Content-Type": contentType || "application/json",
      },
    })
  } catch (error) {
    console.error(`Proxy Error [${request.method}] ${pathname}:`, error)
    return NextResponse.json(
      { code: "ERR_INTERNAL", message: "Failed to connect to backend", details: error instanceof Error ? error.message : String(error) },
      { status: 500 }
    )
  }
}

export const GET = proxyRequest
export const POST = proxyRequest
export const PUT = proxyRequest
export const PATCH = proxyRequest
export const DELETE = proxyRequest
