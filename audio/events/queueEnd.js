const Discord = require("discord.js");
module.exports = {
    name: "queueEnd", // イベント名
    async execute(queue, client) {
        console.log(`The queue has ended.`);
    }
}