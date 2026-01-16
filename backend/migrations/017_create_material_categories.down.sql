DROP TRIGGER IF EXISTS update_material_categories_updated_at ON material_categories;
DROP INDEX IF EXISTS idx_material_categories_deleted_at;
DROP INDEX IF EXISTS idx_material_categories_is_active;
DROP INDEX IF EXISTS idx_material_categories_sort_order;
DROP INDEX IF EXISTS idx_material_categories_code;
DROP TABLE IF EXISTS material_categories;
