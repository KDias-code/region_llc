CREATE TABLE IF NOT EXISTS tasks(
                                    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    id          UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    title       TEXT NOT NULL,
    activeAt    TIMESTAMP
    );