CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

INSERT INTO items (name) VALUES ('Fake Data 1'), ('Fake Data 2');
