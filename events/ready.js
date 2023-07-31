module.exports = {
  name: "ready", // イベント名
  async execute(client) {
    const log = await client.func.readLog("v1/conf/status");
    const add = require(`../system/add.js`);
    add.addCmd(client.conf);
    console.log(`on:${log.onoff},play:${log.playing},status:${log.status}`);
    if (log.onoff == 'on') {
      client.user.setActivity(log.playing);
      client.user.setStatus(log.status);
    } else {
      client.user.setActivity(`Servers: ${client.guilds.cache.size}`);
      client.user.setStatus('online');
    }
    /*client.user.setActivity('activity', { type: 'WATCHING' });
client.user.setActivity('activity', { type: 'LISTENING' });
client.user.setActivity('activity', { type: 'COMPETING' });*/
    /*client.user.setStatus('online');
client.user.setStatus('idle');
client.user.setStatus('dnd');
client.user.setStatus('invisible');*/
    client.channels.cache.get(client.conf.logch.ready).send("Discordログインしました!");
    console.log(`[${client.func.timeToJST(Date.now(), true)}] Logged in as ${client.user.tag}!`);
  }
}