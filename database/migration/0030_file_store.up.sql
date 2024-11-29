-- File Store Table
CREATE TABLE file_store (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    is_system boolean NOT NULL DEFAULT FALSE,
    filename varchar(1000) NOT NULL,
    stored_filename varchar(60) NOT NULL,
    content_type varchar(800),
    file_size bigint NOT NULL,
    browsable boolean NOT NULL DEFAULT TRUE,
    secure boolean NOT NULL DEFAULT FALSE,
    completed_processing boolean NOT NULL DEFAULT FALSE
);

CREATE INDEX file_store_filename_index ON file_store USING HASH (filename);
