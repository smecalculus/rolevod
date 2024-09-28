CREATE TABLE roles (
	id varchar(20) UNIQUE,
	name varchar(20),
	state text
);

CREATE TABLE seats (
	id varchar(20) UNIQUE,
	name varchar(20),
	via text,
	ctx text
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
	on_ref text,
	on_key varchar(20),
	to_id varchar(20),
	to_ids varchar(20)[][2]
);

CREATE TABLE channels (
	id varchar(20),
	name varchar(20),
	pre_id varchar(20),
	state text
);

CREATE TABLE steps (
	id varchar(20) UNIQUE,
	pre_id varchar(20),
	via_id varchar(20),
	kind smallint,
	payload text
);

CREATE TABLE kinships (
	parent_id varchar(20),
	child_id varchar(20)
);
