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

CREATE TABLE crew_roots (
	crew_id varchar(36),
	rev bigint,
	title varchar(64)
);

CREATE TABLE crew_caps (
	crew_id varchar(36),
	sig_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE crew_deps (
	crew_id varchar(36),
	sig_fqn ltree,
	rev_from bigint,
	rev_to bigint
);

CREATE TABLE crew_subs (
	crew_id varchar(36),
	crew_fqn ltree,
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
	sym ltree UNIQUE,
	rev_from bigint,
	rev_to bigint,
	kind smallint
);

CREATE INDEX sym_gist_idx ON aliases USING GIST (sym);
