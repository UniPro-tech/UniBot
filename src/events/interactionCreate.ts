import { Events, EmbedBuilder, Interaction, MessageFlags } from "discord.js";
import config from "@/config";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";

export const name = Events.InteractionCreate;
export const execute = async (interaction: Interaction) => {
  if (interaction.isChatInputCommand()) {
    console.log(
      `[${interaction.client.function.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} info] ->${interaction.commandName}`
    );
    const command = interaction.client.commands.get(interaction.commandName);
    if (!command) {
      console.log(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} info] Not Found: ${interaction.commandName}`
      );
      return;
    }
    if (!interaction.inGuild() && command.guildOnly) {
      const embed = new EmbedBuilder()
        .setTitle("エラー")
        .setDescription("このコマンドはDMでは実行できません。")
        .setColor(interaction.client.config.color.e);
      interaction.reply({ embeds: [embed] });
      console.log(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} info] DM Only: ${interaction.commandName}`
      );
      return;
    }

    try {
      await command.execute(interaction);
      console.log(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} run] ${interaction.commandName}`
      );

      const logEmbed = new EmbedBuilder()
        .setTitle("コマンド実行ログ")
        .setDescription(`${interaction.user} がコマンドを実行しました。`)
        .setColor(config.color.s)
        .setTimestamp()
        .setThumbnail(interaction.user.displayAvatarURL())
        .addFields([
          {
            name: "コマンド",
            value: `\`\`\`\n/${interaction.commandName}\n\`\`\``,
          },
          {
            name: "実行サーバー",
            value:
              "```\n" + interaction.inGuild()
                ? `${interaction.guild?.name} (${interaction.guild?.id})`
                : "DM" + "\n```",
          },
          {
            name: "実行ユーザー",
            value: "```\n" + `${interaction.user.tag}(${interaction.user.id})` + "\n```",
          },
        ])
        .setFooter({ text: `${interaction.id}` });
      const channel = await GetLogChannel(interaction.client);
      if (channel) {
        channel.send({ embeds: [logEmbed] });
      }
    } catch (error) {
      console.error(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} error]An Error Occured in ${interaction.commandName}\nDatails:\n${error}`
      );
      const logEmbed = new EmbedBuilder()
        .setTitle("ERROR - cmd")
        .setDescription("```\n" + (error as any).toString() + "\n```")
        .setColor(config.color.e)
        .setTimestamp();

      const channel = await GetErrorChannel(interaction.client);
      if (channel) {
        channel.send({ embeds: [logEmbed] });
      }
      const messageEmbed = new EmbedBuilder()
        .setTitle("すみません。エラーが発生しました。")
        .setDescription("```\n" + error + "\n```")
        .setColor(config.color.e)
        .setTimestamp();

      await interaction.reply({ embeds: [messageEmbed] });
      const logChannel = await GetLogChannel(interaction.client);
      if (logChannel) {
        logChannel.send({ embeds: [messageEmbed] });
      }
    }
  } else if (interaction.isStringSelectMenu()) {
    if (interaction.customId.startsWith('rp_')) {
      const selected = interaction.values[0];
      console.log(`選択された項目: <@&${selected}>`);
      await interaction.deferReply({ flags: MessageFlags.Ephemeral });
      const member = interaction.guild?.members.cache.get(interaction.user.id);
      if (!member) {
        await interaction.reply({ content: 'メンバー情報を取得できませんでした。', flags: MessageFlags.Ephemeral });
        return;
      }
      const hasRole = member.roles.cache.has(selected);

      if (hasRole) {
        await member.roles.remove(selected);
        await interaction.editReply(`${member.displayName} から役職 <@&${selected}> を削除しました。`);
      } else {
        await member.roles.add(selected);
        await interaction.editReply(`${member.displayName} に役職 <@&${selected}> を付与しました。`);
      }
      console.log(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} info] Selected: ${selected}`
      );
    } else {
      console.log(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} info] Not Found: ${interaction.customId}`
      );
    }
  }
};

export default {
  name,
  execute,
};
