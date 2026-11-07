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
  column "guild_id" {
    type = bigint
    null = false
  }
  column "channel_id" {
    type = bigint
    null = false
  }
  column "title" {
    type = varchar(255)
    null = true
  }
  column "webhook_url" {
    type = varchar(255)
    null = false
  }
  column "url" {
    type = text
    null = false
  }
  column "last_run" {
    type = timestamptz
    null = false
  }
  column "created_at" {
    type    = timestamptz
    null    = false
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_rss_settings_member_id" {
    columns     = [column.user_id, column.guild_id]
    ref_columns = [table.members.column.user_id, table.members.column.guild_id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "fk_rss_settings_channel_id" {
    columns     = [column.channel_id]
    ref_columns = [table.channels.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_last_run" {
    columns = [column.last_run]
  }
}