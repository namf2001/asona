-- ─── ENUMS ───────────────────────────────────────────────────

DO $$ BEGIN
    CREATE TYPE channel_type AS ENUM ('global', 'dm', 'group', 'project');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- ─── CHANNELS ────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.channels (
    id           BIGSERIAL    PRIMARY KEY,
    workplace_id BIGINT                REFERENCES public.workplaces(id) ON DELETE CASCADE,
    project_id   BIGINT                REFERENCES public.projects(id)   ON DELETE CASCADE,
    name         TEXT         CHECK (name <> ''::text),
    type         channel_type NOT NULL,
    created_by   BIGINT                REFERENCES public.users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_channels_workplace ON public.channels(workplace_id);
CREATE INDEX IF NOT EXISTS idx_channels_project ON public.channels(project_id);

-- ─── CHANNEL MEMBERS ─────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.channel_members (
    id           BIGSERIAL   PRIMARY KEY,
    channel_id   BIGINT      NOT NULL REFERENCES public.channels(id) ON DELETE CASCADE,
    user_id      BIGINT      NOT NULL REFERENCES public.users(id)    ON DELETE CASCADE,
    last_read_at TIMESTAMPTZ,                   -- tracks unread count
    joined_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (channel_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_channel_members_user ON public.channel_members(user_id);

-- ─── MESSAGES ────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.messages (
    id         BIGSERIAL   PRIMARY KEY,
    channel_id BIGINT      NOT NULL REFERENCES public.channels(id) ON DELETE CASCADE,
    sender_id  BIGINT               REFERENCES public.users(id)    ON DELETE SET NULL,
    parent_id  BIGINT               REFERENCES public.messages(id) ON DELETE SET NULL,  -- thread reply
    content    TEXT        NOT NULL CHECK (content <> ''::text),
    is_edited  BOOLEAN     NOT NULL DEFAULT false,
    is_deleted BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_channel ON public.messages(channel_id, created_at DESC);

-- ─── MESSAGE ATTACHMENTS ─────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.message_attachments (
    id         BIGSERIAL   PRIMARY KEY,
    message_id BIGINT      NOT NULL REFERENCES public.messages(id) ON DELETE CASCADE,
    file_url   TEXT        NOT NULL CHECK (file_url <> ''::text),
    file_name  TEXT        CHECK (file_name <> ''::text),
    file_size  BIGINT,
    mime_type  TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── MESSAGE REACTIONS ───────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.message_reactions (
    id         BIGSERIAL   PRIMARY KEY,
    message_id BIGINT      NOT NULL REFERENCES public.messages(id) ON DELETE CASCADE,
    user_id    BIGINT      NOT NULL REFERENCES public.users(id)    ON DELETE CASCADE,
    emoji      TEXT        NOT NULL CHECK (emoji <> ''::text),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (message_id, user_id, emoji)
);
