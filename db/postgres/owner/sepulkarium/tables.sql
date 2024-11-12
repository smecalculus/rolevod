CREATE TABLE roles (
	id varchar(36),
	rev bigint,
	name varchar(64),
	state_id varchar(36),
	whole_id varchar(36)
);

CREATE TABLE signatures (
	id varchar(36),
	name varchar(64),
	pe jsonb,
	ces jsonb
);

CREATE TABLE agents (
	id varchar(36),
	name varchar(64)
);

CREATE TABLE deals (
	id varchar(36),
	name varchar(64)
);

CREATE TABLE states (
	id varchar(36),
	kind smallint,
	from_id varchar(36),
	spec jsonb
);

CREATE TABLE channels (
	id varchar(36),
	name varchar(64),
	pre_id varchar(36),
	state_id varchar(36)
);

CREATE TABLE steps (
	id varchar(36),
	kind smallint,
	pid varchar(36),
	vid varchar(36),
	spec jsonb
);

CREATE TABLE kinships (
	parent_id varchar(36),
	child_id varchar(36)
);

CREATE TABLE clientships (
	from_id varchar(36),
	to_id varchar(36),
	pid varchar(36)
);

CREATE TABLE aliases (
	sym ltree UNIQUE,
	id varchar(36),
	rev bigint,
	kind smallint
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
