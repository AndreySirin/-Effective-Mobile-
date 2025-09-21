-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscription (
    subscriptionId UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    serviceName VARCHAR(30) NOT NULL CHECK (LENGTH(serviceName) >= 2),
    price INT NOT NULL,
    userID UUID NOT NULL,
    startDate DATE NOT NULL DEFAULT CURRENT_DATE,
    endDate DATE,
    CHECK (endDate IS NULL OR endDate > startDate)
    );

CREATE INDEX IF NOT EXISTS idx_subscription_user_service_dates_notnull
    ON subscription (userID, serviceName, startDate, endDate)
    INCLUDE (price)
    WHERE endDate IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_subscription_user_service_dates_null
    ON subscription (userID, serviceName, startDate)
    INCLUDE (price)
    WHERE endDate IS NULL;


-- +migrate Down

DROP INDEX IF EXISTS idx_subscription_user_service_dates_null;
DROP INDEX IF EXISTS idx_subscription_user_service_dates_notnull;
DROP TABLE IF EXISTS subscription;
