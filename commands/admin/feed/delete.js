const {
  SlashCommandSubcommandBuilder,
  StringSelectMenuBuilder,
  StringSelectMenuOptionBuilder,
  ActionRowBuilder,
  ComponentType,
} = require("discord.js");
const logUtils = require("../../../lib/logUtils.js");

module.exports = {
  data: new SlashCommandSubcommandBuilder()
    .setName("delete")
    .setDescription("delete"),
  adminGuildOnly: true,
  /**
   * Executes the feed command.
   * @param {CommandInteraction} interaction - The interaction object.
   * @returns {Promise<string>} - A promise that resolves when the execution is complete.
   * @async
   */
  async execute(interaction) {
    const subscribed = await logUtils.readLog(
      "v1/feed/" + interaction.channel.id
    );
    if (!subscribed) {
      interaction.editReply({
        content: "This channel is not subscribed to any feeds.",
        ephemeral: true,
      });
      return "No data";
    }
    const select = new StringSelectMenuBuilder().setCustomId("FeedSelector");
    for (const feed of subscribed.data) {
      await select.addOptions(
        new StringSelectMenuOptionBuilder()
          .setLabel(feed.url)
          .setValue(feed.url)
          .setDescription("Last Update: " + feed.lastDate)
      );
    }

    const row = new ActionRowBuilder().addComponents(select);

    const response = await interaction.editReply({
      content: "Choose your starter!",
      components: [row],
    });

    const collector = response.createMessageComponentCollector({
      componentType: ComponentType.StringSelect,
      time: 3_600_000,
    });

    collector.on("collect", async (i) => {
      const selection = i.values[0];
      const index = subscribed.data.findIndex((x) => x.url === selection);
      subscribed.data.splice(index, 1);
      await logUtils.loging(subscribed, `v1/feed/${interaction.channel.id}`);
      await i.update({
        content: "Deleted: " + selection,
        components: [],
      });
    });
    return "No data";
  },
};
