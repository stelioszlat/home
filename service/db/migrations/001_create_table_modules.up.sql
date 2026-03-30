CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(255),
    isEnabled BOOLEAN DEFAULT true,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_modules_name ON modules(name);

INSERT INTO modules(name, description) VALUES('Analytics', 'Track user behavior and app usage');
INSERT INTO modules(name, description) VALUES('Advanced Logging', 'Enhanced logging capabilities');
INSERT INTO modules(name, description) VALUES('System Monitoring', 'Monitor system health and performance');
INSERT INTO modules(name, description) VALUES('Notifications', 'Advanced system wide notifications');