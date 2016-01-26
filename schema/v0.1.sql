-- INSERT INTO client(id,secret,extra,redirect_uri) VALUES('demo', 'demo', 'b4c8dd5b-852c-460a-9b4a-26109f9162a2', 'http://192.168.1.36:65002/api/v1/oauth2/callback');

CREATE TABLE IF NOT EXISTS clients (
	client_id uuid NOT NULL PRIMARY KEY,
	
	domain text NOT NULL,
	ip4 inet,
	ip6 inet,
	
	secret text NOT NULL,
	redirect text NOT NULL,
	scopes text[] NOT NULL,
	
	flags text[],
	props jsonb NOT NULL DEFAULT '{}',
	
	is_enabled boolean DEFAULT false,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	removed_at timestamp with time zone,
	CONSTRAINT uniq_domains_idx UNIQUE (domain)
);

INSERT INTO clients(client_id, domain, ip4, ip6, secret, redirect, scopes, flags, props, is_enabled) VALUES
	('b4c8dd5b-852c-460a-9b4a-26109f9162a2', 'http://localhost:8081', '127.0.0.1', null, 'demo', 'http://localhost:8081/api/v1/hey/callback', '{demo}', '{demo}', '{}', true);

CREATE TABLE IF NOT EXISTS users (
	user_id uuid PRIMARY KEY,
	client_id uuid,

	ext_id text NOT NULL,
	ext_id_hash text NOT NULL,
	ext_flags text[],
	ext_props jsonb NOT NULL DEFAULT '{}',

	flags text[],
	props jsonb NOT NULL DEFAULT '{}',

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	removed_at timestamp with time zone,
	CONSTRAINT uniq_client_user_idx UNIQUE (client_id, ext_id_hash)
);

CREATE TABLE IF NOT EXISTS channels (
	channel_id uuid PRIMARY KEY,
	client_id uuid,
	
	ext_id text NOT NULL,
	ext_id_hash text NOT NULL,
	ext_flags text[],
	ext_props jsonb NOT NULL DEFAULT '{}',
	
	owners uuid[],
		
	flags text[],
	props jsonb NOT NULL DEFAULT '{}',

	root_thread_id uuid,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	removed_at timestamp with time zone,

	CONSTRAINT uniq_client_channels_idx UNIQUE (client_id, ext_id_hash)
);

CREATE TABLE IF NOT EXISTS channel_counters (
	client_id uuid,
	channel_id uuid,

	counter_events int8,

	CONSTRAINT uniq_client_channel_counter_idx UNIQUE (client_id, channel_id)
);

CREATE TABLE IF NOT EXISTS channel_watchers (
	client_id uuid,
	channel_id uuid,
	user_id uuid,
	
	unread int4,

	CONSTRAINT uniq_client_channel_watchers_idx UNIQUE (client_id, channel_id, user_id)
);	

CREATE TABLE IF NOT EXISTS threads (
	thread_id uuid PRIMARY KEY,
	client_id uuid,
	channel_id uuid,
	
	ext_id text NOT NULL,
	ext_id_hash text NOT NULL,
	ext_flags text[],
	ext_props jsonb NOT NULL DEFAULT '{}',
	
	flags text[],
	props jsonb NOT NULL DEFAULT '{}',

	owners uuid[],

	related_event_id uuid,
	parent_thread_id uuid,
	depth int,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	removed_at timestamp with time zone,

	CONSTRAINT uniq_client_threads_idx UNIQUE (client_id, ext_id_hash)
);

CREATE TABLE IF NOT EXISTS thread_counters (
	client_id uuid,
	thread_id uuid,

	counter_events int8,

	CONSTRAINT uniq_client_thread_counter_idx UNIQUE (client_id, thread_id)
);

CREATE TABLE IF NOT EXISTS thread_watchers (
	client_id uuid,
	thread_id uuid,
	user_id uuid,
	
	unread int4,

	CONSTRAINT uniq_client_thread_watchers_idx UNIQUE (client_id, thread_id, user_id)
);	


CREATE TABLE IF NOT EXISTS threadline (
	client_id uuid,
	channel_id uuid,
	thread_id uuid,

	event_id uuid,

  	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	removed_at timestamp with time zone,
	
  	CONSTRAINT uniq_threads_idx UNIQUE (client_id, channel_id, thread_id, event_id)
);

CREATE INDEX threadline_created_index ON threadline(created_at DESC NULLS LAST);
-- WITH CLUSTERING ORDER BY (created_at DESC, event_id ASC);

CREATE TABLE IF NOT EXISTS events (
	event_id uuid PRIMARY KEY,
	client_id uuid,
	channel_id uuid,
	thread_id uuid,

  	creator uuid,
  	
  	data bytea,

  	props jsonb NOT NULL DEFAULT '{}',

  	parent_thread_id uuid,
  	parent_event_id uuid,
  	branch_thread_id uuid,
	
  	flags text[],
  	ext_flags text[],

  	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	removed_at timestamp with time zone,
	CONSTRAINT uniq_client_event_idx UNIQUE (client_id, event_id)
);