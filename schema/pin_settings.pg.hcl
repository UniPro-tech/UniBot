table "pin_settings" {
  schema = schema.public
  column "channel_id" {
    type = bigint
    null = false
  }
  column "user_id" {
    type = bigint
    null = false
  }
  column "content" {
    type = text
    null = false
  }
  primary_key {
    columns = [column.channel_id]
  }
  foreign_key "fk_channel_id" {
    columns     = [column.channel_id]
    ref_columns = [table.channels.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "fk_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}