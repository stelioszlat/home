CREATE TABLE nodes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    os VARCHAR(100),
    nodeType varchar(100),
    memory varchar(100),
    cpu varchar(100),
    storage varchar(100),
    network varchar(100),
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_nodes_name ON nodes(name);