import { NextRequest, NextResponse } from "next/server";

import { getAuthCookieOptions, getOnboardedCookieOptions } from "@/lib/auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

/**
 * GET /api/auth/callback
 *
 * Handles the Google OAuth redirect from the Go backend.
 * The backend no longer passes the JWT token in the URL (security risk).
 * Instead it passes a short-lived one-time authorization code that is stored
 * in Redis with a 60-second TTL.
 *
 * This route exchanges that code for the real JWT session token via a
 * server-side POST request to Go backend — the token is never visible in any URL.
 */
export async function GET(request: NextRequest) {
  const code = request.nextUrl.searchParams.get("code");

  if (!code) {
    return NextResponse.redirect(new URL("/login?error=missing_code", request.url));
  }

  try {
    // Exchange the one-time code for the real session token.
    // This is a server-side request — token never appears in a URL.
    const res = await fetch(`${API_URL}/api/v1/auth/exchange`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code }),
    });

    if (!res.ok) {
      return NextResponse.redirect(new URL("/login?error=exchange_failed", request.url));
    }

    const body = await res.json();
    const token: string | undefined = body?.data?.session_token;
    const isOnboarded: boolean = body?.data?.is_onboarded === true;

    if (!token) {
      return NextResponse.redirect(new URL("/login?error=invalid_token", request.url));
    }

    const redirectTo = isOnboarded ? "/" : "/onboard";
    const response = NextResponse.redirect(new URL(redirectTo, request.url));

    response.cookies.set("auth_token", token, getAuthCookieOptions());
    response.cookies.set(
      "onboarded",
      isOnboarded ? "true" : "false",
      getOnboardedCookieOptions()
    );

    return response;
  } catch (err) {
    console.error("[OAuth Callback] exchange failed:", err);
    return NextResponse.redirect(new URL("/login?error=server_error", request.url));
  }
}
