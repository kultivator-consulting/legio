-- File Tag Table
CREATE TABLE file_tag (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    file_store_id uuid NOT NULL,
    tag text NOT NULL,
    CONSTRAINT file_tag_file_store_id_fkey FOREIGN KEY (file_store_id) REFERENCES file_store
        ON DELETE CASCADE
);

CREATE INDEX file_tag_file_store_id_index ON file_tag USING HASH (file_store_id);
