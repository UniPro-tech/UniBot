import {
  ChatInputCommandInteraction,
  SlashCommandSubcommandBuilder,
  GuildMember,
  MessageFlags,
  ActionRowBuilder,
  StringSelectMenuBuilder,
  StringSelectMenuOptionBuilder,
  ActionRow,
  StringSelectMenuComponent,
} from "discord.js";
import config from "@/config";
import { readSelected, SelectedDataType } from "@/lib/dataUtils";

export const data = new SlashCommandSubcommandBuilder()
  .setName("add")
  .setDescription("リアクションパネルにロールを追加します")
  .addRoleOption((option) =>
    option.setName("role0").setDescription("追加するロールを選択してください").setRequired(true)
  )
  .addRoleOption((option) =>
    option.setName("role1").setDescription("追加するロールを選択してください(任意)")
  )
  .addRoleOption((option) =>
    option.setName("role2").setDescription("追加するロールを選択してください(任意)")
  )
  .addRoleOption((option) =>
    option.setName("role3").setDescription("追加するロールを選択してください(任意)")
  )
  .addRoleOption((option) =>
    option.setName("role4").setDescription("追加するロールを選択してください(任意)")
  )
  .addRoleOption((option) =>
    option.setName("role5").setDescription("追加するロールを選択してください(任意)")
  )
  .addRoleOption((option) =>
    option.setName("role6").setDescription("追加するロールを選択してください(任意)")
  );

export const adminGuildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const member = interaction.member as GuildMember;
  if (!member.roles.cache.has(config.adminRoleId)) {
    await interaction.reply({ content: "権限がありません", flags: MessageFlags.Ephemeral });
    return; // アドミンロールが付与されていなかったら終了
  }

  const selectedMenuMessageId = (
    await readSelected(interaction.user.id, SelectedDataType.Message)
  )?.data.replace(/"/g, "");
  if (!selectedMenuMessageId) {
    await interaction.reply({
      content: "DB not found",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }
  const selectedMenuMessage = await interaction.channel?.messages.fetch(
    selectedMenuMessageId as string
  );
  if (!selectedMenuMessage) {
    await interaction.reply({
      content: "It is not a valid messageId",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }
  if (!selectedMenuMessage.components || selectedMenuMessage.components.length === 0) {
    await interaction.reply({
      content: "This message does not have any components.",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }
  const selectMenu = (selectedMenuMessage.components[0] as ActionRow<StringSelectMenuComponent>)
    .components[0] as StringSelectMenuComponent;
  if (!selectMenu || selectMenu.type !== 3 || !selectMenu.customId.startsWith("rp_")) {
    await interaction.reply({
      content: "This is not a valid role panel select menu.",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }

  const currentOptions = selectMenu.options;
  const roles = [];
  const memberRoles = member.roles.cache.map((role) => role.position);
  const highestMemberRole = Math.max(...memberRoles);

  const botMember = interaction.guild?.members.me;
  if (!botMember) {
    await interaction.reply({
      content: "ボットのメンバー情報が取得できません",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }
  const botRoles = botMember.roles.cache.map((role) => role.position);
  const highestBotRole = Math.max(...botRoles);

  for (let i = 0; i <= 2; i++) {
    const role = interaction.options.getRole(`role${i}`);
    if (role) {
      // @everyone ロールのIDを取得
      if (role.id === interaction.guild?.id) {
        await interaction.reply({
          content: "`@everyone` ロールは選択できません。",
          flags: MessageFlags.Ephemeral,
        });
        return;
      }

      // ユーザーの役職よりも高い権限のロールを指定した場合
      if (role.position > highestMemberRole) {
        await interaction.reply({
          content: `指定されたロール ${role.name} はあなたより高い権限を持っています。これを付与することはできません。`,
          flags: MessageFlags.Ephemeral,
        });
        return;
      }

      // ボットの役職よりも高い権限のロールを指定した場合
      if (role.position >= highestBotRole) {
        await interaction.reply({
          content: `指定されたロール ${role.name} はこのボットより高い権限を持っています。これを付与することはできません。`,
          flags: MessageFlags.Ephemeral,
        });
        return;
      }

      // BOTロールを指定した場合
      if (role.managed) {
        await interaction.reply({
          content: `指定されたロール ${role.name} はボットロールです。これを付与することはできません。`,
          flags: MessageFlags.Ephemeral,
        });
        return;
      }

      roles.push({
        id: role.id,
        name: role.name,
      });
    }
  }

  for (const option of currentOptions) {
    // 既に選択肢に存在するロールはスキップ
    if (roles.some((role) => role.id === option.value)) {
      continue;
    }

    // 選択肢のロールを追加
    roles.push({
      id: option.value,
      name: option.label,
    });
  }

  // 役職がなければ終了
  if (roles.length === 0) {
    await interaction.reply({
      content: "有効な役職が選択されていません",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }

  // 役職の重複をチェック
  const roleIds = roles.map((role) => role.id);
  const uniqueRoleIds = new Set(roleIds);
  if (uniqueRoleIds.size !== roleIds.length) {
    await interaction.reply({
      content: "同じ役職が複数選択されています",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }

  const panelId = selectMenu.customId;
  const addedSelectMenu = new StringSelectMenuBuilder()
    .setCustomId(`${panelId}`)
    .setPlaceholder("ロールを選択してください")
    .setMinValues(0)
    .setMaxValues(roles.length);

  // 選択肢を追加
  roles.forEach((role) => {
    addedSelectMenu.addOptions(
      new StringSelectMenuOptionBuilder()
        .setLabel(role.name)
        .setValue(role.id.toString())
        .setDescription(`${role.name} ロールを取得/解除します`)
    );
  });

  const row = new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(addedSelectMenu);

  const editEmbed = selectedMenuMessage.embeds[0];

  selectedMenuMessage.edit({
    embeds: [editEmbed],
    components: [row],
  });

  await interaction.reply({
    content: "役職が追加されました",
    flags: MessageFlags.Ephemeral,
  });
};

export default {
  data,
  adminGuildOnly,
};
