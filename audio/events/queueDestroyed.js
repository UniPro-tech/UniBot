const Discord = require("discord.js");
module.exports = {
    name: "queueDestroyed", // イベント名
    async execute(queue, client) {
        console.log(`The queue has destroyed.`);
    }
}