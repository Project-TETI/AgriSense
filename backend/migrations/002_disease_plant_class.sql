BEGIN;

ALTER TABLE diseases ADD COLUMN IF NOT EXISTS symptoms text;
CREATE UNIQUE INDEX IF NOT EXISTS idx_diseases_plant_class ON diseases (plant_id, class_name);
DROP INDEX IF EXISTS idx_diseases_class_name;

COMMIT;
