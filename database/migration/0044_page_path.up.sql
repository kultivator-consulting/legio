-- Page Path Table
CREATE TABLE page_path (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    domain_id uuid NOT NULL, -- the domain to which this page belongs
    account_id uuid NOT NULL, -- the user who created this page path
    parent_page_path_id uuid DEFAULT NULL, -- the parent page path (if any)
    title character varying(120) NOT NULL, -- the name of the page path, e.g. 'Destinations'
    slug character varying(120) NOT NULL, -- the slug of the page path, e.g. 'destinations'
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT page_domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE NO ACTION,
    CONSTRAINT page_page_path_id_fkey FOREIGN KEY (parent_page_path_id) REFERENCES page_path
        ON DELETE NO ACTION,
    CONSTRAINT page_account_id_fkey FOREIGN KEY (account_id) REFERENCES account
        ON DELETE CASCADE
);

CREATE INDEX page_path_title_index ON page_path USING HASH (title);
CREATE INDEX page_path_slug_index ON page_path USING HASH (slug);
CREATE INDEX page_path_domain_id_index ON page_path USING HASH (domain_id);
CREATE INDEX page_path_parent_page_path_id_index ON page_path USING HASH (parent_page_path_id);
CREATE INDEX page_path_account_id_index ON page_path USING HASH (account_id);
