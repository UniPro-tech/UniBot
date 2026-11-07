enum "discord_status_type" {
  schema = schema.public
  values = [
    "online",
    "dnd",
    "idle",
    "invisible",
    "offline",
  ]
}

enum "discord_activity_type" {
  schema = schema.public
  values = [
    "custom",
    "game",
    "competing",
    "listening",
    "watching",
    "streaming",
  ]
}

table "system_preference" {
  schema = schema.public
  column "id" {
    type    = uuid
    null    = false
    default = sql("uuidv7()")
  }
  column "status_type" {
    type    = enum.discord_status_type
    null    = false
    default = "online"
  }
  column "activity_type" {
    type    = enum.discord_activity_type
    null    = true
    default = null
  }
  column "activity_summary" {
    type    = varchar(100)
    null    = true
    default = null
  }
  primary_key {
    columns = [column.id]
  }
}