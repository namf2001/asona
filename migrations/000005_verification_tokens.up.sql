-- ─── VERIFICATION TOKENS ─────────────────────────────────────

CREATE TABLE IF NOT EXISTS public.verification_token (
    identifier TEXT        NOT NULL CHECK (identifier <> ''::text),
    expires    TIMESTAMPTZ NOT NULL,
    token      TEXT        NOT NULL CHECK (token <> ''::text),
    PRIMARY KEY (identifier, token)
);
