-- ─── ENUMS ───────────────────────────────────────────────────

DO $$ BEGIN
    CREATE TYPE project_access AS ENUM ('public', 'private');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE project_role AS ENUM ('owner', 'member');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE task_priority AS ENUM ('highest', 'high', 'medium', 'low', 'lowest');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE task_link_type AS ENUM ('blocks', 'is_blocked_by', 'duplicates', 'relates_to');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- ─── PROJECTS ────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.projects (
    id           BIGSERIAL      PRIMARY KEY,
    workplace_id BIGINT         NOT NULL REFERENCES public.workplaces(id) ON DELETE CASCADE,
    name         TEXT           NOT NULL CHECK (name <> ''::text),
    description  TEXT,
    color        TEXT           CHECK (color <> ''::text),
    access       project_access NOT NULL DEFAULT 'private',
    created_by   BIGINT         REFERENCES public.users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_projects_workplace ON public.projects(workplace_id);

-- ─── PROJECT MEMBERS ─────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.project_members (
    id         BIGSERIAL    PRIMARY KEY,
    project_id BIGINT       NOT NULL REFERENCES public.projects(id) ON DELETE CASCADE,
    user_id    BIGINT       NOT NULL REFERENCES public.users(id)    ON DELETE CASCADE,
    role       project_role NOT NULL DEFAULT 'member',
    joined_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    UNIQUE (project_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_project_members_user ON public.project_members(user_id);

-- ─── PROJECT STATUSES (Kanban columns) ───────────────────────

CREATE TABLE IF NOT EXISTS public.project_statuses (
    id         BIGSERIAL   PRIMARY KEY,
    project_id BIGINT      NOT NULL REFERENCES public.projects(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL CHECK (name <> ''::text),
    color      TEXT        CHECK (color <> ''::text),
    position   INT         NOT NULL DEFAULT 0,
    is_default BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, name)
);

-- ─── TASKS ───────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.tasks (
    id           BIGSERIAL     PRIMARY KEY,
    project_id   BIGINT        NOT NULL  REFERENCES public.projects(id)         ON DELETE CASCADE,
    status_id    BIGINT        NOT NULL  REFERENCES public.project_statuses(id) ON DELETE RESTRICT,
    parent_id    BIGINT                  REFERENCES public.tasks(id)            ON DELETE CASCADE,  -- Subtask
    title        TEXT          NOT NULL  CHECK (title <> ''::text),
    description  TEXT,
    priority     task_priority,
    position     FLOAT         NOT NULL  DEFAULT 0,     -- midpoint-insertion
    due_date     TIMESTAMPTZ,
    start_date   TIMESTAMPTZ,
    estimate     INT,                                   -- estimated minutes/hours
    created_by   BIGINT                  REFERENCES public.users(id) ON DELETE SET NULL,
    reporter_id  BIGINT                  REFERENCES public.users(id) ON DELETE SET NULL,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ   NOT NULL  DEFAULT now(),
    updated_at   TIMESTAMPTZ   NOT NULL  DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tasks_project_status ON public.tasks(project_id, status_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent ON public.tasks(parent_id);

-- ─── TASK ASSIGNEES ──────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.task_assignees (
    task_id     BIGINT      NOT NULL REFERENCES public.tasks(id) ON DELETE CASCADE,
    user_id     BIGINT      NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    assigned_by BIGINT               REFERENCES public.users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, user_id)
);

-- ─── LABELS ──────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.labels (
    id         BIGSERIAL   PRIMARY KEY,
    project_id BIGINT      NOT NULL REFERENCES public.projects(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL CHECK (name <> ''::text),
    color      TEXT        CHECK (color <> ''::text),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, name)
);

CREATE TABLE IF NOT EXISTS public.task_labels (
    task_id  BIGINT NOT NULL REFERENCES public.tasks(id)  ON DELETE CASCADE,
    label_id BIGINT NOT NULL REFERENCES public.labels(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);

-- ─── TASK COMMENTS ───────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.task_comments (
    id         BIGSERIAL   PRIMARY KEY,
    task_id    BIGINT      NOT NULL REFERENCES public.tasks(id)         ON DELETE CASCADE,
    author_id  BIGINT               REFERENCES public.users(id)         ON DELETE SET NULL,
    parent_id  BIGINT               REFERENCES public.task_comments(id) ON DELETE CASCADE,
    content    TEXT        NOT NULL CHECK (content <> ''::text),
    is_edited  BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── TASK ATTACHMENTS ────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.task_attachments (
    id          BIGSERIAL   PRIMARY KEY,
    task_id     BIGINT      NOT NULL REFERENCES public.tasks(id) ON DELETE CASCADE,
    uploaded_by BIGINT               REFERENCES public.users(id) ON DELETE SET NULL,
    file_url    TEXT        NOT NULL CHECK (file_url <> ''::text),
    file_name   TEXT        CHECK (file_name <> ''::text),
    file_size   BIGINT,
    mime_type   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─── TASK LINKS ──────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.task_links (
    id         BIGSERIAL      PRIMARY KEY,
    source_id  BIGINT         NOT NULL REFERENCES public.tasks(id) ON DELETE CASCADE,
    target_id  BIGINT         NOT NULL REFERENCES public.tasks(id) ON DELETE CASCADE,
    link_type  task_link_type NOT NULL,
    created_by BIGINT                  REFERENCES public.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT now(),
    UNIQUE (source_id, target_id, link_type),
    CHECK (source_id <> target_id)
);

-- ─── TASK WATCHERS ───────────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.task_watchers (
    task_id    BIGINT      NOT NULL REFERENCES public.tasks(id) ON DELETE CASCADE,
    user_id    BIGINT      NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, user_id)
);

-- ─── TASK ACTIVITY LOGS ──────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.task_activity_logs (
    id         BIGSERIAL   PRIMARY KEY,
    task_id    BIGINT      NOT NULL REFERENCES public.tasks(id) ON DELETE CASCADE,
    actor_id   BIGINT               REFERENCES public.users(id) ON DELETE SET NULL,
    action     TEXT        NOT NULL   CHECK (action <> ''::text),
    old_value  JSONB,
    new_value  JSONB,
    created_at TIMESTAMPTZ NOT NULL   DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_task_activity_task ON public.task_activity_logs(task_id, created_at DESC);
