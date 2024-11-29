-- File Meta Table
CREATE TABLE file_meta (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    file_store_id uuid NOT NULL,
    key varchar(120) NOT NULL,
    value text DEFAULT NULL,
    attached_file_store_id uuid DEFAULT NULL,
    CONSTRAINT file_meta_file_store_id_fkey FOREIGN KEY (file_store_id) REFERENCES file_store
        ON DELETE CASCADE,
    CONSTRAINT file_meta_attached_file_store_id_fkey FOREIGN KEY (attached_file_store_id) REFERENCES file_store
        ON DELETE CASCADE
);

CREATE INDEX file_meta_file_store_id_index ON file_meta USING HASH (file_store_id);
CREATE INDEX file_meta_attached_file_store_id_index ON file_meta USING HASH (attached_file_store_id);
