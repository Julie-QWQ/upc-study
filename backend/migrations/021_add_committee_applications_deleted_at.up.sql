-- Add missing columns to committee_applications

ALTER TABLE committee_applications
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP;
