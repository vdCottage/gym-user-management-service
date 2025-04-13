-- +migrate Up
-- Add age column as nullable first
ALTER TABLE customers ADD COLUMN age INTEGER;

-- Update age based on date_of_birth for existing records
UPDATE customers 
SET age = EXTRACT(YEAR FROM AGE(CURRENT_DATE, date_of_birth::date))
WHERE date_of_birth IS NOT NULL;

-- Set default age for any remaining NULL values
UPDATE customers SET age = 25 WHERE age IS NULL;

-- Make age column NOT NULL
ALTER TABLE customers ALTER COLUMN age SET NOT NULL;

-- Add profile_url column
ALTER TABLE customers ADD COLUMN profile_url TEXT;

-- Convert health_conditions and fitness_goals to TEXT type
ALTER TABLE customers ALTER COLUMN health_conditions TYPE TEXT;
ALTER TABLE customers ALTER COLUMN fitness_goals TYPE TEXT;

-- Drop date_of_birth column as it's no longer needed
ALTER TABLE customers DROP COLUMN date_of_birth;

-- +migrate Down
-- Add back date_of_birth column
ALTER TABLE customers ADD COLUMN date_of_birth TIMESTAMP;

-- Set approximate date_of_birth based on age
UPDATE customers 
SET date_of_birth = CURRENT_DATE - (age * INTERVAL '1 year');

-- Drop the new columns
ALTER TABLE customers DROP COLUMN age;
-- Make date_of_birth NOT NULL
ALTER TABLE customers ALTER COLUMN date_of_birth SET NOT NULL;

-- Convert back health_conditions and fitness_goals if needed
ALTER TABLE customers 
ALTER COLUMN health_conditions TYPE JSONB USING health_conditions::JSONB,
ALTER COLUMN fitness_goals TYPE TEXT[] USING string_to_array(fitness_goals, ',');

-- Drop new columns
ALTER TABLE customers DROP COLUMN IF EXISTS age;
ALTER TABLE customers DROP COLUMN IF EXISTS profile_url; 