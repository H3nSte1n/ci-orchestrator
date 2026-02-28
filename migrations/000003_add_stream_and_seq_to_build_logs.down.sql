-- Drop indexes
DROP INDEX IF EXISTS idx_build_logs_build_id_created_at;
DROP INDEX IF EXISTS idx_build_logs_build_id_seq;

-- Remove sequence and stream columns
ALTER TABLE build_logs DROP COLUMN seq;
ALTER TABLE build_logs DROP COLUMN stream;
