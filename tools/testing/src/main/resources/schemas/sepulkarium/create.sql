CREATE TABLE sepulkas (
	internal_id UUID UNIQUE,
	external_id CHARACTER VARYING (64) UNIQUE,
	revision INTEGER,
	created_at TIMESTAMP (6),
	updated_at TIMESTAMP (6)
);
