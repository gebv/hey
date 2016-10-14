package postgres

var SchemaBase = []string{
	`CREATE TABLE IF NOT EXISTS channels (
        channel_id uuid PRIMARY KEY,
        client_id uuid,

        ext_id text NOT NULL,
        
        owners uuid[],

        root_thread_id uuid,

        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone DEFAULT now() NOT NULL,

        CONSTRAINT uniq_client_channels_idx UNIQUE (client_id, ext_id)
    )`,
	`CREATE INDEX IF NOT EXISTS uniq_client_channels_idx ON channels(client_id, ext_id) WHERE (ext_id IS NOT NULL)`,
	`CREATE TABLE IF NOT EXISTS threads (
        thread_id uuid PRIMARY KEY,
        client_id uuid,
        channel_id uuid,

        ext_id text NOT NULL,
        
        owners uuid[],

        related_event_id uuid, -- в случае root = nil, в других случая отражает event с которым связан "вверх" поток
	    parent_thread_id uuid, -- в случае root = nil

        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone DEFAULT now() NOT NULL,
        
        CONSTRAINT uniq_client_threads_ids_idx UNIQUE (client_id, thread_id) 
    )`,
	`CREATE INDEX IF NOT EXISTS uniq_client_channel_threads_ext_ids_idx ON threads(client_id, channel_id, ext_id) WHERE (ext_id IS NOT NULL)`,
	`CREATE TABLE IF NOT EXISTS events (
        event_id uuid PRIMARY KEY,
        client_id uuid,
        thread_id uuid,
        channel_id uuid,

        creator uuid,

        data bytea,

        parent_thread_id uuid,
        parent_event_id uuid,
        branch_thread_id uuid,

        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone DEFAULT now() NOT NULL,
        
        CONSTRAINT uniq_client_event_idx UNIQUE (client_id, event_id)
    )`,
	`CREATE TABLE IF NOT EXISTS threadline (
        client_id uuid,
        channel_id uuid,
        thread_id uuid,

        event_id uuid PRIMARY KEY,

        created_at timestamp with time zone NOT NULL
    )`,
	`CREATE INDEX IF NOT EXISTS threadline_created_index ON threadline(
        client_id asc, 
--        channel_id asc, 
        thread_id asc, 
        created_at DESC, 
        event_id ASC
        )`,
	// `CREATE TABLE IF NOT EXISTS thread_watchers (
	//     client_id uuid,
	//     thread_id uuid,
	//     user_id uuid,

	//     unread int8,

	//     CONSTRAINT uniq_client_thread_watchers_idx UNIQUE (client_id, thread_id, user_id)
	// )`,
	// `CREATE TABLE IF NOT EXISTS thread_counters (
	//     client_id uuid,
	//     thread_id uuid,

	//     counter_events int8,

	//     CONSTRAINT uniq_client_thread_counter_idx UNIQUE (client_id, thread_id)
	// )`,
	// `CREATE TABLE IF NOT EXISTS channel_watchers (
	//     client_id uuid,
	//     channel_id uuid,
	//     user_id uuid,

	//     unread int8,

	//     CONSTRAINT uniq_client_channel_watchers_idx UNIQUE (client_id, channel_id, user_id)
	// )`,
	// `CREATE TABLE IF NOT EXISTS channel_counters (
	//     client_id uuid,
	//     channel_id uuid,

	//     counter_events int8,

	//     CONSTRAINT uniq_client_channel_counter_idx UNIQUE (client_id, channel_id)
	// )`,
}
