-- +goose Up
CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    mobile VARCHAR(50),
    home VARCHAR(50)
);

-- +goose Down
DROP TABLE contacts;
