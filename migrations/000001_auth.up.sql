-- ─── USERS ───────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.users (
    id             BIGSERIAL    PRIMARY KEY,
    name           TEXT         CHECK (name <> ''::text),
    username       VARCHAR(100) UNIQUE CHECK (username <> ''::text),
    display_name   VARCHAR(255),
    email          TEXT         UNIQUE CHECK (email <> ''::text),
    email_verified TIMESTAMPTZ,
    password       TEXT,
    avatar_url     TEXT,
    is_active      BOOLEAN      NOT NULL DEFAULT true,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON public.users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON public.users(username);

-- ─── AUTH PROVIDERS ──────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.auth_providers (
    id                  BIGSERIAL   PRIMARY KEY,
    user_id             BIGINT      NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    provider            TEXT        NOT NULL CHECK (provider <> ''::text),
    provider_account_id TEXT        NOT NULL CHECK (provider_account_id <> ''::text),
    access_token        TEXT,
    refresh_token       TEXT,
    token_expires_at    TIMESTAMPTZ,
    id_token            TEXT,
    scope               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_account_id)
);

CREATE INDEX IF NOT EXISTS idx_auth_providers_user ON public.auth_providers(user_id);

-- ─── SESSIONS ────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.sessions (
    id            BIGSERIAL   PRIMARY KEY,
    user_id       BIGINT      NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    session_token TEXT        NOT NULL UNIQUE CHECK (session_token <> ''::text),
    expires_at    TIMESTAMPTZ NOT NULL,
    user_agent    TEXT,
    ip_address    INET,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON public.sessions(user_id);
