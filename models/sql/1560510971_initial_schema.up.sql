BEGIN;

CREATE TABLE issues (
	kallax_id serial NOT NULL PRIMARY KEY,
	id bigint,
	number bigint,
	state text,
	locked boolean,
	title text,
	body text,
	comments bigint,
	closed_at timestamptz,
	created_at timestamptz,
	updated_at timestamptz,
	htmlurl text,
	node_id text,
	repository_owner text NOT NULL,
	repository_name text NOT NULL,
	labels text[] NOT NULL,
	user_id bigint NOT NULL,
	user_login text NOT NULL,
	assignee_id bigint NOT NULL,
	assignee_login text NOT NULL,
	assignees jsonb NOT NULL,
	closed_by_id bigint NOT NULL,
	closed_by_login text NOT NULL,
	milestone_id bigint NOT NULL,
	milestone_title text NOT NULL
);


CREATE TABLE issue_comments (
	kallax_id serial NOT NULL PRIMARY KEY,
	id bigint,
	node_id text,
	body text,
	reactions jsonb,
	created_at timestamptz,
	updated_at timestamptz,
	author_association text,
	htmlurl text,
	user_id bigint NOT NULL,
	user_login text NOT NULL,
	issue_number bigint NOT NULL,
	repository_owner text NOT NULL,
	repository_name text NOT NULL
);


CREATE TABLE organizations (
	kallax_id serial NOT NULL PRIMARY KEY,
	login text,
	id bigint,
	node_id text,
	avatar_url text,
	htmlurl text,
	name text,
	company text,
	blog text,
	location text,
	email text,
	description text,
	public_repos bigint,
	public_gists bigint,
	followers bigint,
	following bigint,
	created_at timestamptz,
	updated_at timestamptz,
	total_private_repos bigint,
	owned_private_repos bigint,
	private_gists bigint,
	disk_usage bigint,
	collaborators bigint,
	billing_email text,
	type text,
	two_factor_requirement_enabled boolean
);


CREATE TABLE pull_requests (
	kallax_id serial NOT NULL PRIMARY KEY,
	id bigint,
	number bigint,
	state text,
	title text,
	body text,
	created_at timestamptz,
	updated_at timestamptz,
	closed_at timestamptz,
	merged_at timestamptz,
	draft boolean,
	merged boolean,
	mergeable boolean,
	mergeable_state text,
	merge_commit_sha text,
	comments bigint,
	commits bigint,
	additions bigint,
	deletions bigint,
	changed_files bigint,
	htmlurl text,
	review_comments bigint,
	maintainer_can_modify boolean,
	author_association text,
	node_id text,
	repository_owner text NOT NULL,
	repository_name text NOT NULL,
	labels text[] NOT NULL,
	user_id bigint NOT NULL,
	user_login text NOT NULL,
	merged_by_id bigint NOT NULL,
	merged_by_login text NOT NULL,
	assignee_id bigint NOT NULL,
	assignee_login text NOT NULL,
	assignees jsonb NOT NULL,
	requested_reviewers jsonb NOT NULL,
	milestone_id bigint NOT NULL,
	milestone_title text NOT NULL,
	head_sha text NOT NULL,
	head_ref text NOT NULL,
	head_label text NOT NULL,
	head_user text NOT NULL,
	head_repository_owner text NOT NULL,
	head_repository_name text NOT NULL,
	base_sha text NOT NULL,
	base_ref text NOT NULL,
	base_label text NOT NULL,
	base_user text NOT NULL,
	base_repository_owner text NOT NULL,
	base_repository_name text NOT NULL
);


CREATE TABLE pull_request_comments (
	kallax_id serial NOT NULL PRIMARY KEY,
	id bigint,
	node_id text,
	in_reply_to bigint,
	body text,
	path text,
	diff_hunk text,
	pull_request_review_id bigint,
	position bigint,
	original_position bigint,
	commit_id text,
	original_commit_id text,
	reactions jsonb,
	created_at timestamptz,
	updated_at timestamptz,
	author_association text,
	htmlurl text,
	user_id bigint NOT NULL,
	user_login text NOT NULL,
	pull_request_number bigint NOT NULL,
	repository_owner text NOT NULL,
	repository_name text NOT NULL
);


CREATE TABLE pull_request_reviews (
	kallax_id serial NOT NULL PRIMARY KEY,
	id bigint,
	node_id text,
	body text,
	submitted_at timestamptz,
	commit_id text,
	htmlurl text,
	state text,
	user_id bigint NOT NULL,
	user_login text NOT NULL,
	pull_request_number bigint NOT NULL,
	repository_owner text NOT NULL,
	repository_name text NOT NULL
);


CREATE TABLE repositories (
	kallax_id serial NOT NULL PRIMARY KEY,
	id bigint,
	node_id text,
	name text,
	full_name text,
	description text,
	homepage text,
	code_of_conduct jsonb,
	default_branch text,
	master_branch text,
	created_at timestamptz,
	pushed_at timestamptz,
	updated_at timestamptz,
	htmlurl text,
	clone_url text,
	git_url text,
	mirror_url text,
	sshurl text,
	svnurl text,
	language text,
	fork boolean,
	forks_count bigint,
	network_count bigint,
	open_issues_count bigint,
	stargazers_count bigint,
	subscribers_count bigint,
	watchers_count bigint,
	size bigint,
	auto_init boolean,
	permissions jsonb,
	allow_rebase_merge boolean,
	allow_squash_merge boolean,
	allow_merge_commit boolean,
	topics text[] NOT NULL,
	archived boolean,
	disabled boolean,
	license jsonb,
	private boolean,
	has_issues boolean,
	has_wiki boolean,
	has_pages boolean,
	has_projects boolean,
	has_downloads boolean,
	license_template text,
	gitignore_template text,
	team_id bigint,
	parent jsonb,
	source jsonb,
	owner_id bigint NOT NULL,
	owner_type text NOT NULL,
	owner_login text NOT NULL,
	organization_id bigint NOT NULL,
	organization_name text NOT NULL
);


CREATE TABLE users (
	kallax_id serial NOT NULL PRIMARY KEY,
	login text,
	id bigint,
	node_id text,
	avatar_url text,
	htmlurl text,
	gravatar_id text,
	name text,
	company text,
	blog text,
	location text,
	email text,
	hireable boolean,
	bio text,
	public_repos bigint,
	public_gists bigint,
	followers bigint,
	following bigint,
	created_at timestamptz,
	updated_at timestamptz,
	suspended_at timestamptz,
	type text,
	site_admin boolean,
	total_private_repos bigint,
	owned_private_repos bigint,
	private_gists bigint,
	disk_usage bigint,
	collaborators bigint,
	two_factor_authentication boolean
);


COMMIT;
