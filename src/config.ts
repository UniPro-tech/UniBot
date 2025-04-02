//(require('dotenv')).config();
const config = {
  adminRoleId: process.env.ROLE_ADMIN,
  color: {
    s: 0x000000,
    e: 0xffffff,
  },
  token: process.env.TOKEN,
  dev: {
    testGuild: process.env.TEST_GUILD,
  },
  logch: {
    ready: process.env.LOG_SUCCESS_ID,
    command: process.env.LOG_SUCCESS_ID,
    error: process.env.LOG_ERROR_ID,
    guildCreate: process.env.LOGCH_GUILD_ID,
    guildDelete: process.env.LOGCH_GUILD_ID,
  },
};

export default config;
