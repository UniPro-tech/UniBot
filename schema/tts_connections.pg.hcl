table "tts_connections" {
  schema = schema.public
  column "guild_id" {
    type = bigint
    null = false
  }
  column "channel_id" {
    type = bigint
    null = false
  }
  primary_key {
    columns = [column.guild_id]
  }
  foreign_key "fk_guild_id" {
    columns     = [column.guild_id]
    ref_columns = [table.guilds.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "fk_channel_id" {
    columns     = [column.channel_id]
    ref_columns = [table.channels.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
}