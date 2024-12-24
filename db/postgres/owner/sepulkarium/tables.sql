CREATE TABLE role_roots (
	role_id varchar(36),
	rev bigint,
	title varchar(64)
);

CREATE TABLE role_states (
	role_id varchar(36),
	state_id varchar(36),
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE role_subs (
	role_id varchar(36),
	role_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE sig_roots (
	sig_id varchar(36),
	rev bigint,
	title text
);

CREATE TABLE sig_pes (
	sig_id varchar(36),
	chnl_key varchar(64),
	role_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE sig_ces (
	sig_id varchar(36),
	chnl_key varchar(64),
	role_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE sig_subs (
	sig_id varchar(36),
	sig_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE pool_roots (
	pool_id varchar(36),
	rev bigint,
	title varchar(64),
	sup_id varchar(64)
);

CREATE TABLE pool_caps (
	pool_id varchar(36),
	sig_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE pool_deps (
	pool_id varchar(36),
	sig_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE pool_subs (
	pool_id varchar(36),
	sub_id varchar(36),
	rev_from bigint,
	rev_to bigint
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

CREATE TABLE producers (
	giver_id varchar(36),
	taker_id varchar(36),
	chnl_id varchar(36)
);

CREATE TABLE consumers (
	giver_id varchar(36),
	taker_id varchar(36),
	chnl_id varchar(36)
);

CREATE TABLE clientships (
	from_id varchar(36),
	to_id varchar(36),
	pid varchar(36)
);

CREATE TABLE aliases (
	id varchar(36),
	sym ltree UNIQUE,
	rev_from bigint,
	rev_to bigint,
	kind smallint
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
