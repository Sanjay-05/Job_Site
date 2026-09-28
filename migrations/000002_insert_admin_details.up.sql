INSERT INTO users (
    id,
    user_name,
    password,
    email,
    role,
    phone_no
) VALUES (
    gen_random_uuid(),
    'admin',
    '$2a$10$EcJoRDPQpLCdaVRcUJNfkOVI5T2.tqlXsLE0gHnmsIx9cYaNkaaT2',
    'admin@example.com',
    'admin',
    '9876543210'
)
ON CONFLICT (email) DO NOTHING;