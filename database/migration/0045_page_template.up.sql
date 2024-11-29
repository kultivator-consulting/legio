-- Page Template Table
CREATE TABLE page_template (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    domain_id uuid NOT NULL, -- the domain to which this page template belongs
    account_id uuid NOT NULL, -- the user who created this page template
    content_id uuid DEFAULT NULL, -- the content for this page template
    parent_page_path_id uuid DEFAULT NULL, -- the parent page path (if any)
    title character varying(120) NOT NULL, -- the name of the page template
    slug character varying(120) NOT NULL, -- the slug of the page template, e.g. 'index'
    description character varying(255) NOT NULL, -- the description of the page template
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT page_template_domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE RESTRICT,
    CONSTRAINT page_template_content_id_fkey FOREIGN KEY (content_id) REFERENCES content
        ON DELETE CASCADE,
    CONSTRAINT page_template_page_path_id_fkey FOREIGN KEY (parent_page_path_id) REFERENCES page_path
        ON DELETE CASCADE,
    CONSTRAINT page_template_account_id_fkey FOREIGN KEY (account_id) REFERENCES account
        ON DELETE NO ACTION
);

CREATE INDEX page_template_title_index ON page_template USING HASH (title);
CREATE INDEX page_template_slug_index ON page_template USING HASH (slug);
CREATE INDEX page_template_domain_id_index ON page_template USING HASH (domain_id);
CREATE INDEX page_template_content_id_index ON page_template USING HASH (content_id);
CREATE INDEX page_template_parent_page_path_id_index ON page_template USING HASH (parent_page_path_id);
CREATE INDEX page_template_account_id_index ON page_template USING HASH (account_id);
