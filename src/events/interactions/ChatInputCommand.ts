import config from "@/config";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { writeChatInputCommandInteractionLog } from "@/lib/logger";
import { ChatInputCommandInteraction, EmbedBuilder } from "discord.js";

const ChatInputCommandExecute = async (interaction: ChatInputCommandInteraction) => {
  const time = () => interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true);
  const log = (msg: string) => console.log(`[${time()}] ${msg}`);

  log(`info ChatInputCommand->${interaction.commandName}`);

  const command = interaction.client.interactionExecutorsCollections.chatInputCommands.get(
    interaction.commandName
  );
  if (!command) {
    log(`info Not Found: ${interaction.commandName}`);
    return;
  }

  if (!interaction.inGuild() && command.guildOnly) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはDMでは実行できません。")
      .setColor(interaction.client.config.color.error);

    await interaction.reply({ embeds: [embed] });
    log(`info DM Only: ${interaction.commandName}`);
    return;
  }

  try {
    await command.execute(interaction);
    log(`run ${interaction.commandName}`);
    await writeChatInputCommandInteractionLog(interaction);
  } catch (error) {
    log(`error An Error Occured in ${interaction.commandName}\nDetails:\n${error}`);

    const errorMsg = (error as any).toString();
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) await errorChannel.send({ embeds: [logEmbed] });

    const userEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    if (!interaction.replied && !interaction.deferred) {
      await interaction.reply({ embeds: [userEmbed] });
    } else {
      await interaction.followUp({ embeds: [userEmbed] });
    }

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) await logChannel.send({ embeds: [userEmbed] });
  }
};

export default ChatInputCommandExecute;
