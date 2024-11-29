-- Content Table
CREATE TABLE content_collection (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    parent_id uuid NOT NULL, -- the parent content instance of this collection, eg. 'a page'
    content_id uuid NOT NULL, -- the content instance, eg. 'a component on a page'
    ordering int NOT NULL, -- the order of the content instance in the collection
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT content_collection_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES content
        ON DELETE CASCADE,
    CONSTRAINT content_collection_content_id_fkey FOREIGN KEY (content_id) REFERENCES content
        ON DELETE CASCADE
);

CREATE INDEX content_collection_parent_id_index ON content_collection USING HASH (parent_id);
CREATE INDEX content_collection_content_id_index ON content_collection USING HASH (content_id);
