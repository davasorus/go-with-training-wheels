CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    priority INTEGER DEFAULT 0, /* 0=Low, 1=Medium, 2=High */
    position INTEGER,
    done BOOLEAN DEFAULT FALSE
);
