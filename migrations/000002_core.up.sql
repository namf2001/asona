-- ─── ENUMS ───────────────────────────────────────────────────

DO $$ BEGIN
    CREATE TYPE verification_code_type AS ENUM ('email_verification', 'organization_join');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE org_role AS ENUM ('admin', 'sub_admin', 'member');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE workplace_role AS ENUM ('admin', 'member');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE workplace_size AS ENUM ('1-10', '11-50', '51-200', '201-500', '500+');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE friendship_status AS ENUM ('pending', 'accepted', 'blocked');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- ─── VERIFICATION CODES ──────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.verification_codes (
    id          BIGSERIAL              PRIMARY KEY,
    user_id     BIGINT                 REFERENCES public.users(id) ON DELETE CASCADE,
    identifier  TEXT        NOT NULL   CHECK (identifier <> ''::text),
    code        TEXT        NOT NULL   CHECK (code <> ''::text),
    type        verification_code_type NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL   DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_verification_codes_identifier ON public.verification_codes(identifier, type);

-- ─── ORGANIZATIONS ───────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.organizations (
    id          BIGSERIAL   PRIMARY KEY,
    name        TEXT        NOT NULL   CHECK (name <> ''::text),
    slug        TEXT        NOT NULL   UNIQUE CHECK (slug <> ''::text),
    logo_url    TEXT,
    description TEXT,
    created_by  BIGINT                 REFERENCES public.users(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL   DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL   DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.organization_members (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT   NOT NULL REFERENCES public.organizations(id) ON DELETE CASCADE,
    user_id         BIGINT   NOT NULL REFERENCES public.users(id)         ON DELETE CASCADE,
    role            org_role NOT NULL DEFAULT 'member',
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_org_members_user ON public.organization_members(user_id);

-- ─── WORKPLACES ──────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.workplaces (
    id         BIGSERIAL      PRIMARY KEY,
    name       TEXT           NOT NULL CHECK (name <> ''::text),
    icon_url   TEXT,
    size       workplace_size,
    created_by BIGINT         REFERENCES public.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.workplace_members (
    id           BIGSERIAL      PRIMARY KEY,
    workplace_id BIGINT         NOT NULL REFERENCES public.workplaces(id) ON DELETE CASCADE,
    user_id      BIGINT         NOT NULL REFERENCES public.users(id)      ON DELETE CASCADE,
    role         workplace_role NOT NULL DEFAULT 'member',
    joined_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),
    UNIQUE (workplace_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_workplace_members_user ON public.workplace_members(user_id);

-- ─── WORKPLACE INVITATIONS ────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.workplace_invitations (
    id           BIGSERIAL   PRIMARY KEY,
    workplace_id BIGINT      NOT NULL REFERENCES public.workplaces(id) ON DELETE CASCADE,
    invite_token TEXT        NOT NULL UNIQUE CHECK (invite_token <> ''::text),
    created_by   BIGINT      REFERENCES public.users(id) ON DELETE SET NULL,
    max_uses     INT,
    use_count    INT         NOT NULL DEFAULT 0,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_workplace_invitations_token ON public.workplace_invitations(invite_token);

-- ─── FRIENDSHIPS ─────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.friendships (
    id           BIGSERIAL        PRIMARY KEY,
    requester_id BIGINT           NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    receiver_id  BIGINT           NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    status       friendship_status NOT NULL DEFAULT 'pending',
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ      NOT NULL DEFAULT now(),
    UNIQUE (requester_id, receiver_id),
    CHECK (requester_id <> receiver_id)
);

CREATE INDEX IF NOT EXISTS idx_friendships_receiver ON public.friendships(receiver_id);
