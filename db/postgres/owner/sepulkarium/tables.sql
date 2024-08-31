CREATE TABLE envs (
	internal_id uuid UNIQUE,
	external_id character varying (64) UNIQUE,
	revision integer,
	created_at timestamp (6),
	updated_at timestamp (6)
);

CREATE TABLE tps (
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE states (
	id varchar(20) UNIQUE,
	name varchar(20),
	kind smallint
);

CREATE TABLE transitions (
	from_id varchar(20),
	to_id varchar(20),
	msg_id varchar(20),
	label varchar(20)
);
