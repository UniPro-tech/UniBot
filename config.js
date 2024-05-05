module.exports = {
  author: {
    name: process.env.AUTHOR_NAME,
    url: process.env.AUTHOR_URL,
    iconURL: process.env.AUTHOR_ICON_URL,
  },
  description: process.env.DESCRIPTION,
  version: process.env.npm_package_version,
  productname: process.env.PRODUCTNAME,
  clientId: process.env.ID,
  adminRoleId: process.env.ADOMIN,
  color: {
    s: 0x1bff49,
    e: 0xff0000,
  },
  token: process.env.DISCORD_TOKEN,
  dev: {
    testGuild: process.env.ADMIN_GUILD,
  },
  logch: {
    ready: Number(process.env.LOG_READY_ID),
    command: Number(process.env.LOG_COMMAND_ID),
    error: Number(process.env.LOG_ERROR_ID),
    guildCreate: Number(process.env.LOG_GUILD_JOIN_ID),
    guildDelete: Number(process.env.LOG_GUID_LEAVE_ID),
  },
  URI_base: process.env.URI_BASE,
};
