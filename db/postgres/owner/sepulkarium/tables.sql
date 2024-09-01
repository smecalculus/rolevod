CREATE TABLE envs (
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE tps (
	id varchar(20) UNIQUE,
	name varchar(20)
);

CREATE TABLE introductions (
	env_id varchar(20),
	tp_id varchar(20)
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
