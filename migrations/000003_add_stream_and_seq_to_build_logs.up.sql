ALTER TABLE build_logs ADD COLUMN stream VARCHAR(10) NOT NULL DEFAULT 'stdout';

ALTER TABLE build_logs ADD COLUMN seq BIGSERIAL UNIQUE;

CREATE INDEX idx_build_logs_build_id_seq ON build_logs(build_id, seq);

CREATE INDEX idx_build_logs_build_id_created_at ON build_logs(build_id, created_at);
