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
  interaction.editReply({
    components: interaction.message.components,
  });

  let completedRoles = Array<{ roleId: string; action: string }>();
  let completed = false;
  try {
    await Promise.all(
      selected.map(async (value) => {
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
      })
    );
    completed = true;
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
};
