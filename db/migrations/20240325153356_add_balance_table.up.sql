CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS balances (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userId uuid DEFAULT uuid_generate_v4(),
    balance int NOT NULL,
    currency TEXT NOT NULL
);