-- Create material_categories table for dynamic category management
CREATE TABLE IF NOT EXISTS material_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_material_categories_code ON material_categories(code);
CREATE INDEX IF NOT EXISTS idx_material_categories_sort_order ON material_categories(sort_order);
CREATE INDEX IF NOT EXISTS idx_material_categories_is_active ON material_categories(is_active);
CREATE INDEX IF NOT EXISTS idx_material_categories_deleted_at ON material_categories(deleted_at);

INSERT INTO material_categories (code, name, description, icon, sort_order) VALUES
    ('courseware', 'Courseware', 'Lecture slides and courseware', 'document', 1),
    ('textbook', 'Textbook', 'Textbooks and course books', 'book', 2),
    ('reference', 'Reference', 'Reference books and materials', 'reference', 3),
    ('exam_paper', 'Exam Paper', 'Past exam papers', 'exam', 4),
    ('exercise', 'Exercise', 'Practice exercises', 'exercise', 5),
    ('experiment', 'Experiment', 'Lab guides and experiment reports', 'experiment', 6),
    ('note', 'Note', 'Class notes and summaries', 'note', 7),
    ('thesis', 'Thesis', 'Academic papers and theses', 'thesis', 8),
    ('other', 'Other', 'Other materials', 'other', 9)
ON CONFLICT (code) DO NOTHING;

CREATE TRIGGER update_material_categories_updated_at
    BEFORE UPDATE ON material_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
