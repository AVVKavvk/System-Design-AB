CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INT NOT NULL CHECk(age > 0 AND age <120),
    create_at TIMESTAMPTZ DEFAULT NOW() 
);