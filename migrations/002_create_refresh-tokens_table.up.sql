CREATE TABLE t_refresh_tokens (
    id              UUID            NOT NULL,
    auth_id         UUID            NOT NULL,
    token           TEXT            NOT NULL,
    expires_at      TIMESTAMP       NOT NULL,
    create_at       TIMESTAMP       NOT NULL    DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (auth_id) REFERENCES t_users(id)
);