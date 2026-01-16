-- Sync materials table with latest model fields

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_category') THEN
        CREATE TYPE material_category AS ENUM (
            'courseware',
            'exam',
            'experiment',
            'exercise',
            'reference',
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
END $$;

ALTER TABLE materials
    ADD COLUMN IF NOT EXISTS category material_category NOT NULL DEFAULT 'other',
    ADD COLUMN IF NOT EXISTS course_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS file_key VARCHAR(500) UNIQUE,
    ADD COLUMN IF NOT EXISTS mime_type VARCHAR(100),
    ADD COLUMN IF NOT EXISTS favorite_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reviewer_id BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'materials' AND column_name = 'file_url'
    ) THEN
        UPDATE materials
        SET file_key = SUBSTRING(file_url FROM '.*/([^/]+)$')
        WHERE file_key IS NULL OR file_key = '';

        ALTER TABLE materials DROP COLUMN IF EXISTS file_url;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'materials' AND column_name = 'search_vector'
    ) THEN
        ALTER TABLE materials ADD COLUMN search_vector tsvector;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_materials_category ON materials(category);
CREATE INDEX IF NOT EXISTS idx_materials_status ON materials(status);
CREATE INDEX IF NOT EXISTS idx_materials_course_name ON materials(course_name);
CREATE INDEX IF NOT EXISTS idx_materials_file_key ON materials(file_key);
CREATE INDEX IF NOT EXISTS idx_materials_reviewer_id ON materials(reviewer_id);
DROP INDEX IF EXISTS idx_materials_title;

DROP INDEX IF EXISTS idx_materials_search;
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
