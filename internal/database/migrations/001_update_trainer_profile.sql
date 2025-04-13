-- Update trainer_profiles table to remove user_id column
ALTER TABLE trainer_profiles DROP COLUMN IF EXISTS user_id; 