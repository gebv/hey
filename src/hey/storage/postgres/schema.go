package postgres

var schemaBase = []string{
	`CREATE TABLE IF NOT EXISTS channels (
        channel_id uuid PRIMARY KEY,
        client_id uuid,
        
        owners uuid[],

        root_thread_id uuid,

        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone DEFAULT now() NOT NULL
    )`,
	`CREATE TABLE IF NOT EXISTS threads (
        thread_id uuid PRIMARY KEY,
        client_id uuid,
        channel_id uuid,
        
        owners uuid[],

        related_event_id uuid, -- в случае root = nil, в других случая отражает event с которым связан "вверх" поток
	    parent_thread_id uuid, -- в случае root = nil

        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone DEFAULT now() NOT NULL
    )`,
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
        channel_id asc, 
        thread_id asc, 
        created_at DESC, 
        event_id ASC
        )`,
}
