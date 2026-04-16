ALTER TABLE public.users
    DROP CONSTRAINT IF EXISTS chk_users_onboarding_step,
    DROP CONSTRAINT IF EXISTS chk_users_onboarding_status,
    DROP COLUMN IF EXISTS onboarding_step,
    DROP COLUMN IF EXISTS onboarding_status;

