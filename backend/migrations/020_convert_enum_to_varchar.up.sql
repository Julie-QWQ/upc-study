-- Convert material_category enum to VARCHAR and enforce FK to material_categories
ALTER TABLE materials
    ALTER COLUMN category TYPE VARCHAR(50) USING category::text;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'materials' AND constraint_name = 'materials_category_fkey'
    ) THEN
        ALTER TABLE materials
            ADD CONSTRAINT materials_category_fkey
            FOREIGN KEY (category)
            REFERENCES material_categories(code)
            ON DELETE RESTRICT
            ON UPDATE CASCADE;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'materials' AND constraint_name = 'materials_category_check'
    ) THEN
        ALTER TABLE materials
            ADD CONSTRAINT materials_category_check
            CHECK (category IS NOT NULL AND category <> '');
    END IF;
END $$;
