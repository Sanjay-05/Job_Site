CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS job_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO countries (id, country_name) VALUES 
(gen_random_uuid(), 'Australia'),
(gen_random_uuid(), 'Canada'),
(gen_random_uuid(), 'India'),
(gen_random_uuid(), 'United Kingdom'),
(gen_random_uuid(), 'United States')
ON CONFLICT (country_name) DO NOTHING;

INSERT INTO job_types (id, type_name) VALUES 
(gen_random_uuid(), 'Full-Time'),
(gen_random_uuid(), 'Freelance'),
(gen_random_uuid(), 'Contract'),
(gen_random_uuid(), 'TeleCom')
ON CONFLICT (type_name) DO NOTHING;
