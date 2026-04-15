/**
 * Parses a Go-style duration string (e.g., "24h", "1h30m", "90m") into seconds.
 * Supports hours (h), minutes (m), and seconds (s).
 */
function parseDurationToSeconds(duration: string): number {
  let totalSeconds = 0;
  const regex = /(\d+)(h|m|s)/g;
  let match;

  while ((match = regex.exec(duration)) !== null) {
    const value = parseInt(match[1], 10);
    const unit = match[2];

    switch (unit) {
      case "h":
        totalSeconds += value * 3600;
        break;
      case "m":
        totalSeconds += value * 60;
        break;
      case "s":
        totalSeconds += value;
        break;
    }
  }

  return totalSeconds;
}

/**
 * Returns the JWT access token max age in seconds, read from the
 * JWT_ACCESS_DURATION env var. Falls back to 24 hours if not set.
 */
export function getTokenMaxAge(): number {
  const duration = process.env.JWT_ACCESS_DURATION || "24h";
  const seconds = parseDurationToSeconds(duration);
  return seconds > 0 ? seconds : 86400; // fallback 24h
}

/**
 * Standard cookie options for the auth token.
 * Centralizes config so every set-cookie call stays consistent.
 */
export function getAuthCookieOptions() {
  return {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax" as const,
    path: "/",
    maxAge: getTokenMaxAge(),
  };
}
