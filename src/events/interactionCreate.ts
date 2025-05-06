import { Events, EmbedBuilder, Interaction, MessageFlags } from "discord.js";
import config from "@/config";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";

export const name = Events.InteractionCreate;
export const execute = async (interaction: Interaction) => {
  if (interaction.isChatInputCommand()) {
    console.log(
      `[${interaction.client.function.timeUtils.timeToJSTstamp(Date.now(), true)} info] ->${
        interaction.commandName
      }`
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
        `[${interaction.client.function.timeUtils.timeToJSTstamp(Date.now(), true)} run] ${
          interaction.commandName
        }`
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
    try {
      if (interaction.customId.startsWith("rp_")) {
        const selected = interaction.values;
        console.log(
          `[${interaction.client.function.timeUtils.timeToJSTstamp(
            Date.now(),
            true
          )} info] -> Menu Selected: ${selected}`
        );
        await interaction.deferUpdate();
        const member = interaction.guild?.members.cache.get(interaction.user.id);
        if (!member) {
          await interaction.followUp({
            content: "メンバー情報を取得できませんでした。",
            flags: MessageFlags.Ephemeral,
          });
          return;
        }
        interaction.editReply({
          components: interaction.message.components,
        });

        let completedRoles = Array<{ roleId: string; action: string }>();
        let completed = false;
        try {
          selected.forEach(async (value, index) => {
            const hasRole = member.roles.cache.has(value);

            if (hasRole) {
              await member.roles.remove(value);
              completedRoles.push({ roleId: value, action: "removed" });
              console.log(
                `[${interaction.client.function.timeUtils.timeToJSTstamp(
                  Date.now(),
                  true
                )} info] -> Role Removed: for ${member.displayName}`
              );
            } else {
              await member.roles.add(value);
              completedRoles.push({ roleId: value, action: "added" });
              console.log(
                `[${interaction.client.function.timeUtils.timeToJSTstamp(
                  Date.now(),
                  true
                )} info] -> Role Added: for ${member.displayName}`
              );
            }
            if (index === selected.length - 1) {
              completed = true;
            }
          });
        } catch (error) {
          console.error(
            `[${interaction.client.function.timeUtils.timeToJSTstamp(
              Date.now(),
              true
            )} error]An Error Occured in ${interaction.customId}\nDetails:\n${error}`
          );
          const messageEmbed = new EmbedBuilder()
            .setTitle("すみません。エラーが発生しました。")
            .setDescription("```\n" + error + "\n```")
            .setColor(config.color.e)
            .setTimestamp();
          await interaction.followUp({
            embeds: [messageEmbed],
            flags: MessageFlags.Ephemeral,
          });
          if (completedRoles.length > 0) {
            const completedRolesString = completedRoles
              .map((role) => {
                return `- <@&${role.roleId}> を ${role.action == "added" ? "追加" : "削除"}`;
              })
              .join("\n");
            await interaction.followUp({
              content: `## 次のとおり変更が完了しました。
<@${interaction.user.id}>さんのロールから
${completedRolesString}
と変更しました。`,
              flags: MessageFlags.Ephemeral,
            });
          }
          return;
        }
        while (!completed) {
          await new Promise((resolve) => setTimeout(resolve, 100));
        }
        if (completedRoles.length > 0) {
          const completedRolesString = completedRoles
            .map((role) => {
              return `- <@&${role.roleId}> を ${role.action == "added" ? "追加" : "削除"}`;
            })
            .join("\n");
          await interaction.followUp({
            content: `## 次のとおり変更が完了しました。
<@${interaction.user.id}>さんのロールから
${completedRolesString}
と変更しました。`,
            flags: MessageFlags.Ephemeral,
          });
        }
      } else {
        console.log(
          `[${interaction.client.function.timeUtils.timeToJSTstamp(
            Date.now(),
            true
          )} info] Not Found: ${interaction.customId}`
        );
      }
    } catch (error) {
      console.error(
        `[${interaction.client.function.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} error]An Error Occured in ${interaction.customId}\nDetails:\n${error}`
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
      if (interaction.channel && interaction.channel.isSendable()) {
        await interaction.channel.send({ embeds: [messageEmbed] });
      }
      const logChannel = await GetLogChannel(interaction.client);
      if (logChannel) {
        logChannel.send({ embeds: [messageEmbed] });
      }
    }
  }
};

export default {
  name,
  execute,
};
