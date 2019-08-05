BEGIN;

CREATE TABLE builds (
	id serial NOT NULL PRIMARY KEY,
	author_date timestamptz,
	author_email text NOT NULL,
	author_name text NOT NULL,
	body text NOT NULL,
	branch text NOT NULL,
	build_num bigint NOT NULL,
	build_parameters jsonb NOT NULL,
	build_time_millis bigint,
	build_url text NOT NULL,
	canceled boolean NOT NULL,
	circle_yml jsonb,
	committer_date timestamptz,
	committer_email text NOT NULL,
	committer_name text NOT NULL,
	compare text,
	dont_build text,
	failed boolean,
	feature_flags jsonb NOT NULL,
	infrastructure_fail boolean NOT NULL,
	is_first_green_build boolean NOT NULL,
	job_name text,
	lifecycle text NOT NULL,
	oss boolean NOT NULL,
	outcome text NOT NULL,
	parallel bigint NOT NULL,
	platform text NOT NULL,
	queued_at text NOT NULL,
	reponame text NOT NULL,
	retries bigint[] NOT NULL,
	retry_of bigint,
	start_time timestamptz,
	status text NOT NULL,
	stop_time timestamptz,
	subject text NOT NULL,
	timedout boolean NOT NULL,
	usage_queued_at text NOT NULL,
	username text NOT NULL,
	vcs_revision text NOT NULL,
	vcs_tag text NOT NULL,
	vcsurl text NOT NULL,
	why text NOT NULL,
	previous_build_num bigint NOT NULL,
	previous_successful_build_num bigint NOT NULL,
	pull_request_urls text[] NOT NULL
);


CREATE TABLE outputs (
	id serial NOT NULL PRIMARY KEY,
	type text NOT NULL,
	time timestamptz NOT NULL,
	message text NOT NULL,
	url text NOT NULL,
	build_num bigint NOT NULL,
	username text NOT NULL,
	reponame text NOT NULL
);


CREATE TABLE steps (
	id serial NOT NULL PRIMARY KEY,
	background boolean NOT NULL,
	bash_command text,
	canceled boolean,
	continue text,
	end_time timestamptz,
	exit_code bigint,
	failed boolean,
	has_output boolean NOT NULL,
	_index bigint NOT NULL,
	infrastructure_fail boolean,
	messages text[] NOT NULL,
	name text NOT NULL,
	output_url text NOT NULL,
	parallel boolean NOT NULL,
	run_time_millis bigint NOT NULL,
	start_time timestamptz,
	status text NOT NULL,
	step bigint NOT NULL,
	timedout boolean,
	truncated boolean NOT NULL,
	type text NOT NULL,
	build_num bigint NOT NULL,
	username text NOT NULL,
	reponame text NOT NULL
);


COMMIT;
