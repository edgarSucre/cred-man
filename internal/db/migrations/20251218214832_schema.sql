-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    birthdate DATE,
    country VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE bank_type AS ENUM ('private', 'government');

CREATE TABLE IF NOT EXISTS banks (
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    type bank_type NOT NULL
);

CREATE TYPE credit_type AS ENUM ('auto', 'mortgage', 'commercial');
CREATE TYPE credit_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE IF NOT EXISTS credit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL,
    bank_id UUID NOT NULL,
    min_payment NUMERIC(19,4) NOT NULL,
    max_payment NUMERIC(19,4) NOT NULL,
    term_months SMALLINT NOT NULL,
    credit_type credit_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status credit_status NOT NULL,
    CONSTRAINT FK_client_id FOREIGN KEY (client_id) REFERENCES clients(id),
    CONSTRAINT FK_bank_id FOREIGN KEY (bank_id) REFERENCES banks(id)
);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS credit;
DROP TABLE IF EXISTS banks;
DROP TABLE IF EXISTS clients;

DROP TYPE IF EXISTS credit_status;
DROP TYPE IF EXISTS credit_type;
DROP TYPE IF EXISTS bank_type;

DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
