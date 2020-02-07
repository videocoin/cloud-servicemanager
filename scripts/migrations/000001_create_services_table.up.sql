CREATE TABLE IF NOT EXISTS services(
    id CHAR (36) PRIMARY KEY,
    name VARCHAR (50) NOT NULL UNIQUE,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL
) ENGINE = INNODB;