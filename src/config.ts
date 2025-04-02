//(require('dotenv')).config();
const config = {
  adminRoleId: process.env.ADMIN_ROLE_ID,
  color: {
    s: 0x000000,
    e: 0xffffff,
  },
  token: process.env.TOKEN,
  dev: {
    testGuild: process.env.TEST_GUILD,
  },
  logch: {
    ready: process.env.LOGCH_READY,
    command: process.env.LOGCH_COMMAND,
    error: process.env.LOGCH_ERROR,
    guildCreate: process.env.LOGCH_GUILD_CREATE,
    guildDelete: process.env.LOGCH_GUILD_DELETE,
  },
};

export default config;
