CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO items (id, name, description)
SELECT
    gen_random_uuid(),
    'Item ' || i,
    'Description for item ' || i
FROM generate_series(1, 100) AS i;
