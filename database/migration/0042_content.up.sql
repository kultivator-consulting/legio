-- Content Table
CREATE TABLE content (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    domain_id uuid NOT NULL, -- the domain to which this content belongs
    component_id uuid NOT NULL, -- the component to which this content belongs
    account_id uuid NOT NULL, -- the user who created this content
    title character varying(120) NOT NULL, -- the name of the content instance, e.g. 'Home Page'
    slug character varying(120) NOT NULL, -- the path or slug of the content instance, e.g. 'home'
    data text NOT NULL, -- the content data TODO can we make this a BJSON field?
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT content_domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE NO ACTION,
    CONSTRAINT content_component_id_fkey FOREIGN KEY (component_id) REFERENCES component
        ON DELETE NO ACTION,
    CONSTRAINT content_account_id_fkey FOREIGN KEY (account_id) REFERENCES account
        ON DELETE NO ACTION
);

CREATE INDEX content_title_index ON content USING HASH (title);
CREATE INDEX content_slug_index ON content USING HASH (slug);
CREATE INDEX content_domain_id_index ON content USING HASH (domain_id);
CREATE INDEX content_component_id_index ON content USING HASH (component_id);
CREATE INDEX content_account_id_index ON content USING HASH (account_id);
