CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE,
    phoneNumber VARCHAR(20) UNIQUE,
    name VARCHAR(100) NOT NULL,
    password VARCHAR(255)NOT NULL,
    CONSTRAINT unique_phone UNIQUE (phoneNumber),
    CONSTRAINT unique_email UNIQUE (email)
);

CREATE INDEX idx_name ON users (name);
-- CREATE INDEX idx_phonenumber ON users (phoneNumber);
-- CREATE INDEX idx_email ON users (email);

