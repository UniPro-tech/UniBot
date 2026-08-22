table "remind_settings" {
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
  column "content" {
    type = text
    null = false
  }
  column "last_run" {
    type = timestamptz

  }
  column "next_run" {
    type = timestamptz
    null = false
  }
  column "cron_expr" {
    type = varchar(100)
    null = true
  }
  column "created_at" {
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
    null    = true
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_remind_settings_member_id" {
    columns     = [column.user_id, column.guild_id]
    ref_columns = [table.members.column.user_id, table.members.column.guild_id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "fk_remind_channel_id" {
    columns     = [column.channel_id]
    ref_columns = [table.channels.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_next_run" {
    columns = [column.next_run]
  }
}