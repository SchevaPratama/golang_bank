CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS transactions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userId uuid DEFAULT uuid_generate_v4(),
    bankAccountNumber varchar(255) NOT NULL,
    bankName varchar(255) NOT NULL,
    transferProofImage varchar(255) DEFAULT '',
    balance int NOT NULL,
    currency TEXT NOT NULL,
    flow TEXT NOT NULL,
    createdAt timestamptz NOT NULL default current_timestamp,
    updatedAt timestamptz NOT NULL default current_timestamp
);