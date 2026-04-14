"use server";

import { cookies } from "next/headers";
import { z } from "zod";
import { loginSchema, passwordSchema } from "./schema";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";

export async function loginAction(data: z.infer<typeof loginSchema>) {
  try {
    const res = await fetch(`${API_URL}/api/v1/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: data.email, password: data.password }),
    });

    const body = await res.json();

    if (!res.ok) {
      return { error: body.message || "Login failed" };
    }

    const { session_token } = body.data || {};
    if (session_token) {
      const cookieStore = await cookies();
      cookieStore.set("auth_token", session_token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        path: "/",
      });
    }

    return { success: true };
  } catch (error: unknown) {
    console.error("Login Action Error:", error);
    return { error: "Network or server connection error" };
  }
}

export async function registerStep1Action(email: string) {
  try {
    const res = await fetch(`${API_URL}/api/v1/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ step: 1, email }),
    });

    const body = await res.json();
    if (!res.ok) {
      return { error: body.message || "Failed to send verification code" };
    }

    return { success: true };
  } catch (error) {
    console.error("Register Step 1 Error:", error);
    return { error: "Network or server connection error" };
  }
}

export async function registerStep2Action(email: string, otp: string) {
  try {
    const res = await fetch(`${API_URL}/api/v1/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ step: 2, email, otp }),
    });

    const body = await res.json();
    if (!res.ok) {
      return { error: body.message || "Invalid or expired verification code" };
    }

    return { success: true };
  } catch (error) {
    console.error("Register Step 2 Error:", error);
    return { error: "Network or server connection error" };
  }
}

export async function registerStep3Action(input: {
  email: string;
  otp: string;
  name: string;
  password: z.infer<typeof passwordSchema>;
}) {
  const username = input.email.split("@")[0] + Math.floor(Math.random() * 10000);

  try {
    const res = await fetch(`${API_URL}/api/v1/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        step: 3,
        email: input.email,
        otp: input.otp,
        name: input.password.name,
        username,
        password: input.password.password,
      }),
    });

    const body = await res.json();
    if (!res.ok) {
      return { error: body.message || "Registration failed" };
    }

    const { session_token } = body.data || {};
    if (session_token) {
      const cookieStore = await cookies();
      cookieStore.set("auth_token", session_token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "lax",
        path: "/",
      });
    }

    return { success: true };
  } catch (error) {
    console.error("Register Step 3 Error:", error);
    return { error: "Network or server connection error" };
  }
}

export async function setAuthTokenAction(token: string) {
  const cookieStore = await cookies();
  cookieStore.set("auth_token", token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
  });
  return { success: true };
}
