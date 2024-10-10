CREATE TABLE roles (
	id varchar(20),
	name varchar(64),
	state jsonb
);

CREATE TABLE seats (
	id varchar(20),
	name varchar(64),
	via jsonb,
	ctx jsonb
);

CREATE TABLE agents (
	id varchar(20),
	name varchar(64)
);

CREATE TABLE deals (
	id varchar(20),
	name varchar(64)
);

CREATE TABLE participations (
	id varchar(20),
	deal_id varchar(20),
	seat_id varchar(20),
	pak varchar(36),
	cak varchar(36)
);

CREATE TABLE states (
	id varchar(20),
	kind smallint,
	from_id varchar(20),
	fqn varchar(512),
	pair jsonb,
	choices jsonb
);

CREATE TABLE channels (
	id varchar(20),
	name varchar(64),
	pre_id varchar(20),
	st_id varchar(20),
	state jsonb
);

CREATE TABLE steps (
	id varchar(20),
	kind smallint,
	pid varchar(20),
	vid varchar(20),
	ctx jsonb,
	term jsonb
);

CREATE TABLE kinships (
	parent_id varchar(20),
	child_id varchar(20)
);

CREATE TABLE clientships (
	pid varchar(20),
	to_id varchar(20),
	from_id varchar(20)
);

CREATE TABLE aliases (
	sym ltree UNIQUE,
	id varchar(20)
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
