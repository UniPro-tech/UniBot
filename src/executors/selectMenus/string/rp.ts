import { EmbedBuilder, MessageFlags, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";

export const name = "rp";

export const execute = async (interaction: StringSelectMenuInteraction) => {
  const selected = interaction.values;
  console.log(
    `[${interaction.client.functions.timeUtils.timeToJSTstamp(
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

  await interaction.editReply({
    components: interaction.message.components,
  });

  const completedRoles: { roleId: string; action: string }[] = [];

  try {
    for (const value of selected) {
      const hasRole = member.roles.cache.has(value);
      if (hasRole) {
        await member.roles.remove(value);
        completedRoles.push({ roleId: value, action: "removed" });
        console.log(
          `[${interaction.client.functions.timeUtils.timeToJSTstamp(
            Date.now(),
            true
          )} info] -> Role Removed: for ${member.displayName}`
        );
      } else {
        await member.roles.add(value);
        completedRoles.push({ roleId: value, action: "added" });
        console.log(
          `[${interaction.client.functions.timeUtils.timeToJSTstamp(
            Date.now(),
            true
          )} info] -> Role Added: for ${member.displayName}`
        );
      }
    }
  } catch (error) {
    console.error(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occured in ${interaction.customId}\nDetails:\n${error}`
    );
    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.followUp({
      embeds: [messageEmbed],
      flags: MessageFlags.Ephemeral,
    });
  }

  if (completedRoles.length > 0) {
    const completedRolesString = completedRoles
      .map((role) => `- <@&${role.roleId}> を ${role.action === "added" ? "追加" : "削除"}`)
      .join("\n");
    await interaction.followUp({
      content: `## 次のとおり変更が完了しました。
<@${interaction.user.id}>さんのロールから
${completedRolesString}
と変更しました。`,
      flags: MessageFlags.Ephemeral,
    });
  }
};
