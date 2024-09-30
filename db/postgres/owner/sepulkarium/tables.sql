CREATE TABLE roles (
	id varchar(20) UNIQUE,
	name varchar(20),
	state jsonb
);

CREATE TABLE seats (
	id varchar(20) UNIQUE,
	name varchar(20),
	via jsonb,
	ctx jsonb
);

CREATE TABLE agents (
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE deals (
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE participations (
	part_id varchar(20),
	deal_id varchar(20),
	seat_id varchar(20),
	pak varchar(36),
	cak varchar(36)
);

CREATE TABLE states (
	id varchar(20) UNIQUE,
	kind smallint,
	from_id varchar(20),
	on_ref jsonb,
	to_id varchar(20),
	choices jsonb
);

CREATE TABLE channels (
	id varchar(20),
	name varchar(20),
	pre_id varchar(20),
	state jsonb
);

CREATE TABLE steps (
	id varchar(20) UNIQUE,
	kind smallint,
	via_id varchar(20),
	payload jsonb
);

CREATE TABLE kinships (
	parent_id varchar(20),
	child_id varchar(20)
);

CREATE TABLE aliases (
	sym ltree UNIQUE,
	id varchar(20)
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
