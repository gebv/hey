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
}
