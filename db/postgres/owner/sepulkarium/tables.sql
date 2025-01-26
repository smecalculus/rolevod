CREATE TABLE role_roots (
	role_id varchar(36),
	title varchar(64),
	rev bigint
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
	title text,
	rev bigint
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
	title varchar(64),
	sup_id varchar(36),
	rev bigint,
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

CREATE TABLE pool_acts (
	pool_id varchar(36),
	act_key varchar(36), -- ключ распоряжения, пользования...
	kind smallint,
	spec jsonb,
	rev_at bigint
);

CREATE TABLE chnl_roots (
	chnl_id varchar(36),
	title varchar(64),
	revs integer[]  -- states = 1, pools = 2
);

CREATE TABLE chnl_states (
	chnl_id varchar(36),
	state_id varchar(36),
	act_key varchar(36),
	rev integer
);

CREATE TABLE chnl_pools (
	chnl_id varchar(36),
	pool_id varchar(36),
	rev integer
);

CREATE TABLE deals (
	id varchar(36),
	name varchar(64)
);

CREATE TABLE states (
	id varchar(36),
	from_id varchar(36),
	kind smallint,
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
