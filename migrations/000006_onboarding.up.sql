-- ─── ONBOARDING ───────────────────────────────────────────────
-- Tracks when a user has completed the onboarding flow.
-- NULL means the user hasn't finished onboarding yet.

ALTER TABLE public.users
    ADD COLUMN IF NOT EXISTS onboarded_at TIMESTAMPTZ;
