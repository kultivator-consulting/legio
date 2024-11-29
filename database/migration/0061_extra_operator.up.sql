-- Extra Operator Mapping Table
CREATE TABLE extra_operator (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  extra_id uuid NOT NULL,
                                  operator_id uuid NOT NULL,
                                  CONSTRAINT extra_id_fkey FOREIGN KEY (extra_id) REFERENCES extra
                                      ON DELETE CASCADE,
                                  CONSTRAINT operator_id_fkey FOREIGN KEY (operator_id) REFERENCES operator
                                      ON DELETE CASCADE
);

CREATE INDEX extra_operator_extra_id_index ON extra_operator USING HASH (extra_id);
CREATE INDEX extra_operator_operator_id_index ON extra_operator USING HASH (operator_id);
