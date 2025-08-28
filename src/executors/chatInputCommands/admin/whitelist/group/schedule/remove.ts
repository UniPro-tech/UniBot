import { readConfig, writeConfig } from "@/lib/dataUtils";
import { ChatInputCommandInteraction, SlashCommandSubcommandBuilder } from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("remove")
  .setDescription("ホワイトリストからロールもしくはユーザーを削除")
  .addRoleOption((option) =>
    option.setName("role").setDescription("削除するロール").setRequired(false)
  )
  .addUserOption((option) =>
    option.setName("user").setDescription("削除するユーザー").setRequired(false)
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const config = await readConfig("whitelist:schedule");
  const roles = Array.isArray(config?.roles) ? config.roles : [];
  const users = Array.isArray(config?.users) ? config.users : [];

  let updatedRoles = roles;
  const role = interaction.options.getRole("role");
  if (role) {
    updatedRoles = roles.filter((id: string) => id !== role.id.toString());
  }

  let updatedUsers = users;
  const user = interaction.options.getUser("user");
  if (user) {
    updatedUsers = users.filter((id: string) => id !== user.id.toString());
  }

  await writeConfig(
    {
      roles: updatedRoles,
      users: updatedUsers,
    },
    "whitelist:schedule"
  );
  await interaction.reply({
    content: `ホワイトリストから${role ? `ロール「${role.name}」` : ""}${
      user ? `ユーザー「${user.username}」` : ""
    }を削除しました。`,
    ephemeral: true,
  });
};
