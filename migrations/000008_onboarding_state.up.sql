-- ─── ONBOARDING STATE ─────────────────────────────────────────
-- Adds explicit onboarding state fields for flow orchestration.

ALTER TABLE public.users
    ADD COLUMN IF NOT EXISTS onboarding_status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS onboarding_step SMALLINT NOT NULL DEFAULT 0;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_users_onboarding_status'
    ) THEN
        ALTER TABLE public.users
            ADD CONSTRAINT chk_users_onboarding_status
            CHECK (onboarding_status IN ('pending', 'in_progress', 'completed'));
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_users_onboarding_step'
    ) THEN
        ALTER TABLE public.users
            ADD CONSTRAINT chk_users_onboarding_step
            CHECK (onboarding_step >= 0 AND onboarding_step <= 3);
    END IF;
END $$;

UPDATE public.users
SET onboarding_status = 'completed', onboarding_step = 3
WHERE onboarded_at IS NOT NULL;

