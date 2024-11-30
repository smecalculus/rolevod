CREATE TABLE roles (
	id varchar(36),
	rev bigint,
	name varchar(64),
	state_id varchar(36),
	whole_id varchar(36)
);

CREATE TABLE role_roots (
	role_id varchar(36),
	role_rev bigint,
	role_name varchar(64),
	role_desc varchar(64),
	state_id varchar(36),
	whole_id varchar(36)
);

CREATE TABLE role_snaps (
	role_id varchar(36),
	role_rev bigint,
	role_name varchar(64),
	state_id varchar(36),
	whole_id varchar(36)
);

CREATE TABLE role_states (
	role_id varchar(36),
	rev_from bigint,
	rev_to bigint,
	state_id varchar(36)
);

CREATE TABLE role_fqns (
	role_id varchar(36),
	rev_from bigint,
	rev_to bigint,
	fqn ltree
);

CREATE TABLE signatures (
	id varchar(36),
	name varchar(64),
	pe jsonb,
	ces jsonb
);

CREATE TABLE sig_roots (
	sig_id varchar(36),
	rev bigint,
	title text
);

CREATE TABLE sig_snaps (
	sig_id varchar(36),
	rev bigint,
	pe jsonb,
	ces jsonb
);

CREATE TABLE sig_pes (
	sig_id varchar(36),
	rev_from bigint,
	rev_to bigint,
	chnl_key varchar(64),
	role_fqn ltree
);

CREATE TABLE sig_ces (
	sig_id varchar(36),
	rev_from bigint,
	rev_to bigint,
	chnl_key varchar(64),
	role_fqn ltree
);

CREATE TABLE agents (
	id varchar(36),
	name varchar(64)
);

CREATE TABLE teams (
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
	id varchar(36),
	rev_from bigint,
	rev_to bigint,
	sym ltree UNIQUE,
	kind smallint
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
