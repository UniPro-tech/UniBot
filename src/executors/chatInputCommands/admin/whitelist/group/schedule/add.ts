import { readConfig, writeConfig } from "@/lib/dataUtils";
import { ChatInputCommandInteraction, SlashCommandSubcommandBuilder } from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("add")
  .setDescription("ホワイトリストへロールもしくはユーザーを追加")
  .addRoleOption((option) =>
    option.setName("role").setDescription("追加するロール").setRequired(false)
  )
  .addUserOption((option) =>
    option.setName("user").setDescription("追加するユーザー").setRequired(false)
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const config = await readConfig("whitelist:schedule");
  const roles = Array.isArray(config?.roles) ? config.roles : [];
  const users = Array.isArray(config?.users) ? config.users : [];

  let appendedRoles = roles;
  const role = interaction.options.getRole("role");
  if (role) {
    appendedRoles = Array.from(new Set([...roles, role.id.toString()]));
  }

  let appendedUsers = users;
  const user = interaction.options.getUser("user");
  if (user) {
    appendedUsers = Array.from(new Set([...users, user.id.toString()]));
  }

  await writeConfig(
    {
      roles: appendedRoles,
      users: appendedUsers,
    },
    "whitelist:schedule"
  );
  await interaction.reply({
    content: `ホワイトリストへ${
      interaction.options.getRole("role")
        ? `ロール「${interaction.options.getRole("role")?.name}」`
        : ""
    }${
      interaction.options.getUser("user")
        ? `ユーザー「${interaction.options.getUser("user")?.username}」`
        : ""
    }を追加しました。`,
    ephemeral: true,
  });
};
