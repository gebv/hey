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
}
