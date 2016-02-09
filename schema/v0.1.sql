-- INSERT INTO client(id,secret,extra,redirect_uri) VALUES('demo', 'demo', 'b4c8dd5b-852c-460a-9b4a-26109f9162a2', 'http://192.168.1.36:65002/api/v1/oauth2/callback');

-- CREATE DATABASE hey_access;
-- CREATE DATABASE hey_app;

/*
id           text NOT NULL PRIMARY KEY,
	secret 		 text NOT NULL,
	extra 		 text NOT NULL,
	redirect_uri text NOT NULL
*/
CREATE TABLE IF NOT EXISTS client (
	id uuid NOT NULL PRIMARY KEY,
	
	domain text NOT NULL,
	ips inet[],
	
	secret text NOT NULL,
	redirect_uri text NOT NULL,
	scopes text[] NOT NULL,
	
	flags text[],
	props jsonb NOT NULL DEFAULT '{}',
	
	is_enabled boolean DEFAULT false,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
	CONSTRAINT uniq_domains_idx UNIQUE (domain)
);

INSERT INTO client(id, domain, ips, secret, redirect_uri, scopes, flags, props, is_enabled, created_at) VALUES
	('b4c8dd5b-852c-460a-9b4a-26109f9162a2', 'http://localhost:8081', '{127.0.0.1}', 'demo', 'http://localhost:8081/api/v1/hey/callback', '{demo}', '{demo}', '{}', true, now());

CREATE TABLE IF NOT EXISTS users (
	user_id uuid PRIMARY KEY,
	client_id uuid,

	ext_id text NOT NULL,
	ext_id_hash text NOT NULL,
	ext_props jsonb NOT NULL DEFAULT '{}',

	is_enabled boolean DEFAULT false,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,
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

	root_thread_id uuid,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,

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
	
	unread int8,

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

	related_event_id uuid, -- в случае root = nil, в других случая отражает event с которым связан "вверх" поток
	parent_thread_id uuid, -- в случае root = nil
	depth smallint,

	is_removed boolean DEFAULT false,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone DEFAULT now() NOT NULL,

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
	
	unread int8,

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
	
  	CONSTRAINT uniq_threads_idx UNIQUE (client_id, channel_id, thread_id, event_id)
);

CREATE INDEX threadline_created_index ON threadline(created_at DESC NULLS LAST, event_id ASC);
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

	CONSTRAINT uniq_client_event_idx UNIQUE (client_id, event_id)
);