-- Update trainers table to make gym_owner_id nullable initially
ALTER TABLE trainers ALTER COLUMN gym_owner_id DROP NOT NULL;

-- Update existing trainers to have a default gym owner if needed
-- You may want to set this to an actual gym owner ID in your system
-- UPDATE trainers SET gym_owner_id = 'default_gym_owner_uuid' WHERE gym_owner_id IS NULL;

-- After data migration, if you want to make it required again:
-- ALTER TABLE trainers ALTER COLUMN gym_owner_id SET NOT NULL; 