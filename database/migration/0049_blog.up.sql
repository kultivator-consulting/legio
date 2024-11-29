-- Blog Table
CREATE TABLE blog (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    domain_id uuid NOT NULL, -- the domain to which this blog belongs
    account_id uuid NOT NULL, -- the user who created this blog
    page_id uuid NOT NULL, -- the page to which this blog belongs
    title character varying(120) NOT NULL,
    description text NOT NULL DEFAULT '',
    image varchar(400) NOT NULL,
    image_info varchar(600) NOT NULL,
    keywords text[],
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT blog_domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE NO ACTION,
    CONSTRAINT blog_account_id_fkey FOREIGN KEY (account_id) REFERENCES account
        ON DELETE NO ACTION,
    CONSTRAINT blog_page_id_fkey FOREIGN KEY (page_id) REFERENCES page
        ON DELETE NO ACTION
);

CREATE INDEX blog_title_index ON blog USING HASH (title);
CREATE INDEX blog_domain_id_index ON blog USING HASH (domain_id);
CREATE INDEX blog_account_id_index ON blog USING HASH (account_id);
CREATE INDEX blog_page_id_index ON blog USING HASH (page_id);
