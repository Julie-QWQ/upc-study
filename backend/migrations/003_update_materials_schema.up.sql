-- Update materials schema for extended management features

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_category') THEN
        CREATE TYPE material_category AS ENUM (
            'textbook',
            'reference',
            'exam_paper',
            'note',
            'exercise',
            'experiment',
            'thesis',
            'other'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_status') THEN
        CREATE TYPE material_status AS ENUM (
            'pending',
            'approved',
            'rejected',
            'deleted'
        );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'report_status') THEN
        CREATE TYPE report_status AS ENUM (
            'pending',
            'approved',
            'rejected'
        );
    END IF;
END $$;

ALTER TABLE materials
    ALTER COLUMN category TYPE material_category USING category::text::material_category,
    ALTER COLUMN status TYPE material_status USING status::text::material_status,
    ADD COLUMN IF NOT EXISTS course_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS file_key VARCHAR(500) UNIQUE,
    ADD COLUMN IF NOT EXISTS mime_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS favorite_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reviewer_id BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_materials_file_key ON materials(file_key);
CREATE INDEX IF NOT EXISTS idx_materials_course_name ON materials(course_name);
CREATE INDEX IF NOT EXISTS idx_materials_reviewer_id ON materials(reviewer_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'materials' AND column_name = 'search_vector'
    ) THEN
        ALTER TABLE materials ADD COLUMN search_vector tsvector;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_materials_search ON materials USING gin(search_vector);

CREATE OR REPLACE FUNCTION materials_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('simple', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('simple', COALESCE(NEW.course_name, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS materials_search_vector_trigger ON materials;
CREATE TRIGGER materials_search_vector_trigger
    BEFORE INSERT OR UPDATE ON materials
    FOR EACH ROW
    EXECUTE FUNCTION materials_search_vector_update();

UPDATE materials SET search_vector =
    setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('simple', COALESCE(description, '')), 'B') ||
    setweight(to_tsvector('simple', COALESCE(course_name, '')), 'C')
WHERE search_vector IS NULL;

ALTER TABLE reports
    RENAME COLUMN reporter_id TO user_id;

ALTER TABLE reports
    ALTER COLUMN status DROP DEFAULT;

ALTER TABLE reports
    ALTER COLUMN status TYPE report_status USING status::text::report_status;

ALTER TABLE reports
    ALTER COLUMN status SET DEFAULT 'pending';

ALTER TABLE reports
    ADD COLUMN IF NOT EXISTS handler_id BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS handled_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS handle_note TEXT;

DROP INDEX IF EXISTS idx_reports_reporter_id;
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);
CREATE INDEX IF NOT EXISTS idx_reports_handler_id ON reports(handler_id);

ALTER TABLE favorites
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

ALTER TABLE download_records
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE TRIGGER update_favorites_updated_at
    BEFORE UPDATE ON favorites
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_download_records_updated_at
    BEFORE UPDATE ON download_records
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
