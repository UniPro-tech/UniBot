import {
  CommandInteraction,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
  EmbedBuilder,
} from "discord.js";

import config from "@/config";

export const data = new SlashCommandSubcommandBuilder()
  .setName("add")
  .setDescription("ReactionPanelのロールを追加します")
  .addRoleOption((option) =>
    option.setName("role0").setDescription("付与するロールを選択してください").setRequired(true)
  )
  .addRoleOption((option) =>
    option.setName("role1").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role2").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role3").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role4").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role5").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role6").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role7").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role8").setDescription("付与するロールを選択してください").setRequired(false)
  )
  .addRoleOption((option) =>
    option.setName("role9").setDescription("付与するロールを選択してください").setRequired(false)
  );
export const adminGuildOnly = true;
export const execute = async (interaction: CommandInteraction) => {
  if (!(interaction.member?.roles as GuildMemberRoleManager).cache.has(config.adminRoleId)) {
    interaction.reply("権限がありません");
    return;
  }
  await interaction.reply("作成中");
};

export default {
  data,
  adminGuildOnly,
  execute,
};
