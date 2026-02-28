CREATE TABLE builds
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_url TEXT NOT NULL,
    ref TEXT NOT NULL,
    command TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    exit_code INT DEFAULT 0,
    error TEXT,
    locked_by TEXT,
    locked_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    cancel_requested_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE build_logs
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    build_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (build_id) REFERENCES builds (id) ON DELETE CASCADE
);
