import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
  PermissionsBitField,
  MessageFlags,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";
import { readConfig } from "@/lib/dataUtils";
import { loggingSystem } from "@/index";

export const handlingCommands = subCommandHandling("general/schedule");
export const data = addSubCommand(
  "general/schedule",
  new SlashCommandBuilder().setName("schedule").setDescription("予約投稿を管理します。")
);
export const guildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const logger = loggingSystem.getLogger({ function: "general/schedule" });
  const whitelist = await readConfig("whitelist:schedule");
  const allowedRoles = Array.isArray(whitelist?.roles) ? whitelist.roles : [];
  const allowedUsers = Array.isArray(whitelist?.users) ? whitelist.users : [];
  if (
    !(
      !interaction.guild ||
      (interaction.guild &&
        (() => {
          // member.rolesはGuildMemberRoleManagerまたはstring[]
          const rolesRaw = interaction.member?.roles;
          let roleIds: string[] = [];
          if (
            rolesRaw &&
            typeof rolesRaw === "object" &&
            "cache" in rolesRaw &&
            rolesRaw.cache instanceof Map
          ) {
            // GuildMemberRoleManager
            roleIds = Array.from(rolesRaw.cache.keys());
          } else if (Array.isArray(rolesRaw)) {
            // string[]
            roleIds = rolesRaw;
          }
          return roleIds.some((roleId) => allowedRoles.includes(roleId));
        })()) ||
      allowedUsers.includes(interaction.user.id) ||
      interaction.memberPermissions?.has(PermissionsBitField.Flags.Administrator)
    )
  ) {
    await replyWithError(
      interaction,
      "Error - 権限がありません",
      "このコマンドを実行する権限がありません。"
    );
    return;
  }

  const subcommand = (interaction.options as CommandInteractionOptionResolver).getSubcommand();
  const command = handlingCommands.get(subcommand);

  if (!command) {
    logger.error({ context: { command: interaction.commandName } }, "No command handler found");
    return;
  }

  try {
    await command.execute(interaction);
    logger.info({ context: { command: interaction.commandName } }, "Command executed successfully");
  } catch (error) {
    logger.error(
      { context: { command: interaction.commandName }, stack_trace: (error as Error).stack },
      "Command execution failed",
      error
    );

    const errorMsg = (error as Error).toString();
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) {
      errorChannel.send({ embeds: [logEmbed] });
    }

    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    await interaction.reply({ embeds: [messageEmbed] });

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
};

const replyWithError = async (
  interaction: ChatInputCommandInteraction,
  title: string,
  description: string
) => {
  const embed = new EmbedBuilder()
    .setTitle(title)
    .setDescription(description)
    .setColor(config.color.error)
    .setTimestamp();
  await interaction.reply({
    embeds: [embed],
    flags: [MessageFlags.Ephemeral],
  });
};
