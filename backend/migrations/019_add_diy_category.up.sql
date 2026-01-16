-- Add diy category to material_category enum
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'diy';
