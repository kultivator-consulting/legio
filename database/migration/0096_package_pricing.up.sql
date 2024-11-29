-- Package Pricing Table
CREATE TABLE package_pricing (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    package_id uuid NOT NULL,
    supplement varchar(30) NOT NULL DEFAULT 'Twin/Double', -- ENUM 'Single', 'Twin/Double', 'Triple', 'Quad', 'Child with Bed', 'Child without Bed'
    adult_price decimal NOT NULL,
    child_price decimal DEFAULT 0.00,
    infant_price decimal DEFAULT 0.00,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    specific_terms text NOT NULL,
    CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
        ON DELETE CASCADE
);

CREATE INDEX package_pricing_package_id_index ON package_pricing USING HASH (package_id);
