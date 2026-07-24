table "rss_settings" {
  schema = schema.public
  column "id" {
    type    = uuid
    null    = false
    default = sql("uuidv7()")
  }
  column "user_id" {
    type = bigint
    null = false
  }
  column "channel_id" {
    type = bigint
    null = false
  }
  column "url" {
    type = text
    null = false
  }
  column "last_run" {
    type = timestamptz
    null = true
  }
  column "created_at" {
    type    = timestamptz
    null    = false
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "fk_channel_id" {
    columns     = [column.channel_id]
    ref_columns = [table.channels.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_last_run" {
    columns = [column.last_run]
  }
}