-- Page Table
CREATE TABLE page (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    domain_id uuid NOT NULL, -- the domain to which this page belongs
    account_id uuid NOT NULL, -- the user who created this page
    content_id uuid DEFAULT NULL, -- the content for this page (if any)
    page_path_id uuid DEFAULT NULL, -- the parent page path (if any)
    title character varying(120) NOT NULL, -- the name of the page, e.g. 'Home Page'
    slug character varying(120) NOT NULL, -- the slug of the page, e.g. 'index'
    seo_title character varying(120) NOT NULL, -- the title of the page for SEO purposes
    seo_description character varying(255) NOT NULL, -- the description of the page for SEO purposes
    seo_keywords character varying(255) NOT NULL, -- the keywords of the page for SEO purposes
    draft_page_id uuid DEFAULT NULL, -- the page this page is a draft of (if any)
    page_template_id uuid DEFAULT NULL, -- the page template this page is based on
    publish_at timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'), -- the date and time at which the page should be published
    unpublish_at timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'), -- the date and time at which the page should be unpublished
    version integer NOT NULL DEFAULT 1, -- the version of the page
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT page_domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE RESTRICT,
    CONSTRAINT page_content_id_fkey FOREIGN KEY (content_id) REFERENCES content
        ON DELETE CASCADE,
    CONSTRAINT page_page_path_id_fkey FOREIGN KEY (page_path_id) REFERENCES page_path
        ON DELETE CASCADE,
    CONSTRAINT page_account_id_fkey FOREIGN KEY (account_id) REFERENCES account
        ON DELETE NO ACTION,
    CONSTRAINT page_page_template_id_fkey FOREIGN KEY (page_template_id) REFERENCES page_template
        ON DELETE NO ACTION,
    CONSTRAINT page_page_id_fkey FOREIGN KEY (draft_page_id) REFERENCES page
        ON DELETE CASCADE
);

CREATE INDEX page_title_index ON page USING HASH (title);
CREATE INDEX page_slug_index ON page USING HASH (slug);
CREATE INDEX page_domain_id_index ON page USING HASH (domain_id);
CREATE INDEX page_content_id_index ON page USING HASH (content_id);
CREATE INDEX page_page_path_id_index ON page USING HASH (page_path_id);
CREATE INDEX page_account_id_index ON page USING HASH (account_id);
CREATE INDEX page_page_id_index ON page USING HASH (draft_page_id);
CREATE INDEX page_page_template_id_index ON page USING HASH (page_template_id);
CREATE INDEX page_publish_at_index ON page USING HASH (publish_at);
CREATE INDEX page_unpublish_at_index ON page USING HASH (unpublish_at);
