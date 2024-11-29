-- Page Path Extension Mapping Table
CREATE TABLE page_path_extension (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    page_path_id uuid NOT NULL,
    extension_id uuid NOT NULL,
    CONSTRAINT page_path_id_fkey FOREIGN KEY (page_path_id) REFERENCES page_path
        ON DELETE CASCADE,
    CONSTRAINT extension_id_fkey FOREIGN KEY (extension_id) REFERENCES extension
        ON DELETE NO ACTION
);

CREATE INDEX page_path_extension_page_path_id_index ON page_path_extension USING HASH (page_path_id);
CREATE INDEX page_path_extension_extension_id_index ON page_path_extension USING HASH (extension_id);
