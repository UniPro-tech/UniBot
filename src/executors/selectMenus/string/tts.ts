import { EmbedBuilder, MessageFlags, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";
import { subSelectMenusHandling } from "@/lib/commandUtils";
import { GetErrorChannel } from "@/lib/channelUtils";

export const name = "tts";

const handlingCommands = subSelectMenusHandling("string/tts");

export const execute = async (interaction: StringSelectMenuInteraction) => {
  const commands = handlingCommands.get(interaction.customId.split("_")[1]);
  if (!commands) {
    console.log(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} info] Not Found: ${interaction.customId}`
    );
    return;
  }
  commands.execute(interaction).catch(async (error) => {
    console.error(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred in ${interaction.customId}\nDetails:\n${error}`
    );
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + (error as any).toString() + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    const channel = await GetErrorChannel(interaction.client);
    if (channel) {
      channel.send({ embeds: [logEmbed] });
    }
    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.error)
      .setTimestamp();
    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({
        embeds: [messageEmbed],
        flags: MessageFlags.SuppressEmbeds,
      });
    }
  });
};
