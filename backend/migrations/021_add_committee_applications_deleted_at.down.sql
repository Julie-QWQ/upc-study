-- Remove columns added in 021

ALTER TABLE committee_applications
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS deleted_at;
