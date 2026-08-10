DROP TABLE IF EXISTS todo;
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    XXXXX VARCHAR(128) NOT NULL,
    YYY VARCHAR(255) NOT NULL,
    ZZZ DECIMAL(5, 2) NOT NULL
);
INSERT INTO todos (XXXXX, YYY, ZZZ)
VALUES ('Blue Train', 'John Coltrane', 56.99),
    ('Giant Steps', 'John Coltrane', 63.99),
    ('Jeru', 'Gerry Mulligan', 17.99),
    ('Sarah Vaughan', 'Sarah Vaughan', 34.98);