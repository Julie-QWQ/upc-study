-- Update review workflow related tables

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_status') THEN
        CREATE TYPE application_status AS ENUM (
            'pending',
            'approved',
            'rejected',
            'cancelled'
        );
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
        CREATE TYPE notification_type AS ENUM (
            'system',
            'material',
            'committee',
            'report'
        );
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_status') THEN
        CREATE TYPE notification_status AS ENUM (
            'unread',
            'read'
        );
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'review_action') THEN
        CREATE TYPE review_action AS ENUM (
            'approve',
            'reject'
        );
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'review_target') THEN
        CREATE TYPE review_target AS ENUM (
            'material',
            'committee',
            'report'
        );
    END IF;
END $$;

ALTER TABLE committee_applications
    ALTER COLUMN status DROP DEFAULT;

ALTER TABLE committee_applications
    ALTER COLUMN status TYPE application_status USING status::text::application_status;

ALTER TABLE committee_applications
    ALTER COLUMN status SET DEFAULT 'pending';

ALTER TABLE committee_applications
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP;

ALTER TABLE committee_applications
    DROP CONSTRAINT IF EXISTS committee_applications_user_id_key;
CREATE INDEX IF NOT EXISTS idx_committee_applications_user_id_status
    ON committee_applications(user_id, status);

ALTER TABLE review_records
    ALTER COLUMN action TYPE review_action USING action::text::review_action;

ALTER TABLE review_records
    ADD COLUMN IF NOT EXISTS target_type review_target,
    ADD COLUMN IF NOT EXISTS target_id BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS original_data JSONB,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_review_records_target_type ON review_records(target_type);
CREATE INDEX IF NOT EXISTS idx_review_records_target_id ON review_records(target_id);
DROP INDEX IF EXISTS idx_review_records_material_id;

ALTER TABLE notifications
    ALTER COLUMN type TYPE notification_type USING type::text::notification_type;

ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS status notification_status NOT NULL DEFAULT 'unread',
    ADD COLUMN IF NOT EXISTS link VARCHAR(255),
    ADD COLUMN IF NOT EXISTS read_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
DROP INDEX IF EXISTS idx_notifications_is_read;

ALTER TABLE reports
    ALTER COLUMN status TYPE report_status USING status::text::report_status;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'reports'
          AND column_name = 'handle_comment'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'reports'
          AND column_name = 'handle_note'
    ) THEN
        EXECUTE 'ALTER TABLE reports RENAME COLUMN handle_comment TO handle_note';
    END IF;
END $$;

DROP TRIGGER IF EXISTS update_committee_applications_updated_at ON committee_applications;
CREATE TRIGGER update_committee_applications_updated_at
    BEFORE UPDATE ON committee_applications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_review_records_updated_at ON review_records;
CREATE TRIGGER update_review_records_updated_at
    BEFORE UPDATE ON review_records
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
CREATE TRIGGER update_notifications_updated_at
    BEFORE UPDATE ON notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
