const { SlashCommandBuilder } = require("discord.js");
const Discord = require("discord.js");
const logUtils = require("../../lib/logUtils");
const { rssGet } = require("../../lib/rss.cjs");

module.exports = {
  guildOnly: true,
  adminGuildOnly: true,
  data: new SlashCommandBuilder()
    .setName("feed")
    .setDescription("Regist RSS feed")
    .addStringOption((option) =>
      option.setName("url").setDescription("URL of RSS feed").setRequired(true)
    ),
  /**
   * Executes the feed command.
   * @param {CommandInteraction} i - The interaction object.
   * @returns {Promise<string>} - A promise that resolves when the execution is complete.
   * @async
   */
  async execute(i) {
    if (i.member.permissions.has("ADMINISTRATOR") === false) {
      i.reply("You don't have permission to use this command.");
      return "No permission";
    }
    const feed = await rssGet(i.options.getString("url"));
    const embed = new Discord.EmbedBuilder()
      .setTitle("Registed RSS feed")
      .addFields([
        {
          name: "URL",
          value: ` ** ${i.options.getString("url")} ** `,
          inline: true,
        },
      ])
      .setFields({
        name: "Title",
        value: ` ** ${feed[0].title} ** `,
        inline: true,
      })
      .setFields({
        name: "FirstContent",
        value: ` ** ${feed[0].content} ** `,
        inline: true,
      })
      .setColor(i.client.conf.color.s)
      .setTimestamp();
    i.reply({ embeds: [embed] });
    console.log("OK");
    try {
      let log = await logUtils.readLog("v1/feed/" + i.channel.id);
      if (log) {
        log.data.push({
          url: i.options.getString("url"),
          lastDate: feed[0].pubDate,
        });
        await logUtils.loging(log, `v1/feed/${i.channel.id}`);
      } else {
        log = {
          data: [
            { url: i.options.getString("url"), lastDate: feed[0].pubDate },
          ],
        };
        await logUtils.loging(log, `v1/feed/${i.channel.id}`);
      }
    } catch (e) {
      console.error(e);
    }
    console.log("OK2");
    return "No data";
  },
};
