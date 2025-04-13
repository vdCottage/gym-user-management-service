-- Create trainers table
CREATE TABLE trainers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'trainer',
    is_active BOOLEAN NOT NULL DEFAULT false,
    gym_owner_id UUID REFERENCES gym_owners(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create trainer_profiles table
CREATE TABLE trainer_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trainer_id UUID NOT NULL REFERENCES trainers(id) ON DELETE CASCADE,
    specialization VARCHAR(100) NOT NULL,
    experience INTEGER NOT NULL DEFAULT 0,
    bio TEXT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    rating DECIMAL(3,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_trainers_gym_owner_id ON trainers(gym_owner_id);
CREATE INDEX idx_trainers_email ON trainers(email);
CREATE INDEX idx_trainers_phone ON trainers(phone);
CREATE INDEX idx_trainer_profiles_trainer_id ON trainer_profiles(trainer_id);
CREATE INDEX idx_trainer_profiles_specialization ON trainer_profiles(specialization); 