-- Extend material_category enum values
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'textbook';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'note';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'exercise';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'thesis';
ALTER TYPE material_category ADD VALUE IF NOT EXISTS 'exam_paper';
