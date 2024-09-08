CREATE TABLE roles (
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE seats (
	id varchar(20) UNIQUE,
	name varchar(20)
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
	deal_id varchar(20),
	seat_id varchar(20)
);

CREATE TABLE states (
	kind smallint,
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE transitions (
	from_id varchar(20),
	to_id varchar(20),
	msg_id varchar(20),
	msg_key varchar(20)
);

CREATE TABLE channels (
	id varchar(20),
	name varchar(20),
	state_id varchar(20)
);

CREATE TABLE kinships (
	parent_id varchar(20),
	child_id varchar(20)
);
