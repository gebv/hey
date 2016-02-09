package models

//go:generate goco -type=file_models -config=goco_helpful.json -out=goco_helpful.go

//go:generate goco -type=file_models -config=goco_config.json -out=goco_config.go

//go:generate goco -type=file_models -config=goco_session.json -out=goco_session.go

//go:generate goco -type=file_models -config=goco_context.json -out=goco_context.go


//go:generate goco -type=file_models -config=goco_client.json -out=goco_client.go

//go:generate goco -type=file_models -config=goco_users.json -out=goco_users.go

//go:generate goco -type=file_models -config=goco_channels.json -out=goco_channels.go

//go:generate goco -type=file_models -config=goco_threads.json -out=goco_threads.go

//go:generate goco -type=file_models -config=goco_events.json -out=goco_events.go
