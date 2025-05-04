-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE project_status AS ENUM ('proposed', 'approved', 'invested', 'disbursed');

CREATE TABLE IF NOT EXISTS projects
(
    id                     UUID PRIMARY KEY        DEFAULT uuid_generate_v4(),
    name                   VARCHAR(255)   NOT NULL,

    borrower_id            VARCHAR(10)    NOT NULL,
    borrower_name          VARCHAR(255)   NOT NULL,
    borrower_mail          VARCHAR(255)   NOT NULL,
    borrower_agreement_url TEXT,

    current_status         project_status NOT NULL DEFAULT 'proposed',
    current_pic_name       VARCHAR(255)   NOT NULL,
    current_pic_mail       VARCHAR(255)   NOT NULL,

    loan_principal_amount  NUMERIC(12, 2) NOT NULL,
    total_invested_amount  NUMERIC(12, 2) NOT NULL DEFAULT 0,
    borrower_rate          NUMERIC(5, 2)  NOT NULL,
    roi_rate               NUMERIC(5, 2)  NOT NULL,

    approved_at            TIMESTAMP,
    disbursed_at           TIMESTAMP,

    created_at             TIMESTAMP      NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMP      NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMP      NULL     DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS project_histories
(
    id         UUID PRIMARY KEY                       DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects (id) NOT NULL,
    status     project_status                NOT NULL,
    pic_name   VARCHAR(255),
    pic_mail   VARCHAR(255),
    extra      JSONB                         NOT NULL DEFAULT '{}',

    created_at TIMESTAMP                     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP                     NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP                     NULL     DEFAULT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS unx_project_histories ON project_histories (project_id, status);

CREATE TABLE IF NOT EXISTS project_investments
(
    id                       UUID PRIMARY KEY                       DEFAULT uuid_generate_v4(),
    project_id               UUID REFERENCES projects (id) NOT NULL,

    investor_id              UUID                          NOT NULL,
    investor_name            VARCHAR(255)                  NOT NULL,
    investor_mail            VARCHAR(255)                  NOT NULL,

    investment_amount        NUMERIC(12, 2)                NOT NULL,
    investment_agreement_url TEXT                          NOT NULL,

    created_at               TIMESTAMP                     NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMP                     NOT NULL DEFAULT NOW(),
    deleted_at               TIMESTAMP                     NULL     DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_histories;
DROP TABLE IF EXISTS project_investments;
DROP TABLE IF EXISTS projects;
DROP TYPE IF EXISTS project_status;
-- +goose StatementEnd
