import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { ButtonInteraction, EmbedBuilder } from "discord.js";
import config from "@/config";

const ButtonExecute = async (interaction: ButtonInteraction) => {
  const time = () => interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true);

  try {
    const [prefix] = interaction.customId.split("_");
    const executionDefine = interaction.client.interactionExecutorsCollections.buttons.get(prefix);

    if (!executionDefine) {
      console.log(`[${time()} info] Not Found: ${interaction.customId}`);
      return;
    }

    console.log(`[${time()} info] Button -> ${interaction.customId}`);
    await executionDefine.execute(interaction);
  } catch (error) {
    const errorMsg = (error as Error).toString();
    console.error(
      `[${time()} error] An Error Occured in ${interaction.customId}\nDetails:\n${errorMsg}`
    );

    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const [errorChannel, logChannel] = await Promise.all([
      GetErrorChannel(interaction.client),
      GetLogChannel(interaction.client),
    ]);

    if (errorChannel) errorChannel.send({ embeds: [logEmbed] });
    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({ embeds: [messageEmbed] });
    }
    if (logChannel) logChannel.send({ embeds: [messageEmbed] });
  }
};

export default ButtonExecute;
