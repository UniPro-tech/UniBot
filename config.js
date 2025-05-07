//(require('dotenv')).config();
module.exports = {
  adminRoleId: process.env.ROLE_ADMIN,
  color: {
    s: 0x1bff49,
    e: 0xff0000,
  },
  token: process.env.DISCORD_TOKEN,
  dev: {
    testGuild: process.env.ADMIN_GUILD,
  },
  logch: {
    ready: process.env.LOG_SUCCESS_ID,
    command: process.env.LOG_SUCCESS_ID,
    error: process.env.LOG_ERROR_ID,
    guildCreate: process.env.LOG_GUILD_ID,
    guildDelete: process.env.LOG_GUID_ID,
  },
};
