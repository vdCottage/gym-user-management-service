-- +migrate Up
-- Create gym_owners table
CREATE TABLE gym_owners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    gym_name VARCHAR(100) NOT NULL,
    gym_registration_number VARCHAR(20) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create customers table
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(15) UNIQUE NOT NULL,
    gender VARCHAR(10) NOT NULL,
    age INT NOT NULL,
    height FLOAT NOT NULL,
    weight FLOAT NOT NULL,
    health_conditions TEXT,
    fitness_goals TEXT,
    profile_url TEXT,
    role VARCHAR(50) NOT NULL DEFAULT 'customer',
    is_active BOOLEAN DEFAULT TRUE,
    gym_owner_id UUID NOT NULL REFERENCES gym_owners(id),
    steak INT NOT NULL DEFAULT 0,
    recommended_diet_template JSON DEFAULT NULL,
    recommended_exercise_template JSON DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create trainers table
CREATE TABLE trainers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'trainer',
    is_active BOOLEAN DEFAULT TRUE,
    gym_owner_id UUID NOT NULL REFERENCES gym_owners(id),
    specialization VARCHAR(100) NOT NULL,
    experience INT NOT NULL DEFAULT 0,
    bio TEXT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    rating FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- +migrate Down
-- Drop trainers table
DROP TABLE IF EXISTS trainers;

-- Drop customers table
DROP TABLE IF EXISTS customers;

-- Drop gym_owners table
DROP TABLE IF EXISTS gym_owners;