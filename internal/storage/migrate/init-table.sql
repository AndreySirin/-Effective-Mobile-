CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate up

CREATE TABLE IF NOT EXISTS subscription (
    subscriptionId UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    serviceName VARCHAR(30) NOT NULL CHECK (LENGTH(serviceName) >= 2),
    price INT,
    userID UUID NOT NULL,
    startDate DATE NOT NULL DEFAULT CURRENT_DATE,
    );

-- +migrate down

DROP TABLE IF EXISTS subscription;