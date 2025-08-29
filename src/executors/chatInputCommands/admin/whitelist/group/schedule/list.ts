import { readConfig } from "@/lib/dataUtils";
import {
  ChatInputCommandInteraction,
  MessageFlags,
  SlashCommandSubcommandBuilder,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("list")
  .setDescription("ホワイトリストに登録されているロール・ユーザーの一覧を表示");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const config = await readConfig("whitelist:schedule");
  const roles = Array.isArray(config?.roles) ? config.roles : [];
  const users = Array.isArray(config?.users) ? config.users : [];

  const roleMentions =
    roles.length > 0 ? roles.map((id: string) => `<@&${id}>`).join("\n") : "(なし)";
  const userMentions =
    users.length > 0 ? users.map((id: string) => `<@${id}>`).join("\n") : "(なし)";

  await interaction.reply({
    content: `**ホワイトリスト一覧**\n\n**ロール:**\n${roleMentions}\n\n**ユーザー:**\n${userMentions}`,
    flags: [MessageFlags.Ephemeral],
  });
};
