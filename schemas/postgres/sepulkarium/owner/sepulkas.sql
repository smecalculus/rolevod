CREATE TABLE sepulkas (
	internal_id uuid UNIQUE,
	external_id character varying (64) UNIQUE,
	revision integer,
	created_at timestamp (6),
	updated_at timestamp (6)
);
