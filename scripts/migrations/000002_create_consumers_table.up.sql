CREATE TABLE IF NOT EXISTS consumers(
    id VARCHAR (50) PRIMARY KEY,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL
) ENGINE = INNODB;