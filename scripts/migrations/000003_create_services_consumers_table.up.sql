CREATE TABLE IF NOT EXISTS services_consumers(
    service_id CHAR (36),
    consumer_id VARCHAR (50),
    FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE,
    FOREIGN KEY (consumer_id) REFERENCES consumers (id) ON DELETE CASCADE,
    PRIMARY KEY (service_id, consumer_id)
) ENGINE = INNODB;