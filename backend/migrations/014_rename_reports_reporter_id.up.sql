-- Rename reports.reporter_id to user_id when it exists

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'reports' AND column_name = 'reporter_id'
    ) THEN
        ALTER TABLE reports RENAME COLUMN reporter_id TO user_id;
    END IF;
END $$;

DROP INDEX IF EXISTS idx_reports_reporter_id;
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'reports' AND column_name = 'handled_at'
    ) THEN
        ALTER TABLE reports ADD COLUMN handled_at TIMESTAMP;
    END IF;
END $$;
