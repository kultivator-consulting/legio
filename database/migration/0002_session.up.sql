-- Session Table
CREATE TABLE session (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    account_id uuid DEFAULT NULL,
    client_id uuid NOT NULL,
    client_agent varchar(1600) NOT NULL,
    client_ip varchar(200) NOT NULL,
    client_bundle_id varchar(1600),
    access_token varchar(1600) NOT NULL,
    refresh_token varchar(1600) NOT NULL,
    access_token_expiry timestamp with time zone NOT NULL,
    refresh_token_expiry timestamp with time zone NOT NULL,
    is_device_app boolean NOT NULL DEFAULT false,
    CONSTRAINT session_account_id_fkey FOREIGN KEY (account_id) REFERENCES account
        ON DELETE NO ACTION
);

CREATE INDEX session_account_id_index ON session USING HASH (account_id);
CREATE INDEX session_client_id_index ON session USING HASH (client_id);
CREATE INDEX session_access_token_index ON session USING HASH (access_token);
CREATE INDEX session_refresh_token_index ON session USING HASH (refresh_token);
