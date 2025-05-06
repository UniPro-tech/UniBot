import {
  ChatInputCommandInteraction,
  SlashCommandSubcommandBuilder,
  EmbedBuilder,
  GuildMember,
  MessageFlags,
  ActionRowBuilder,
  StringSelectMenuBuilder,
  StringSelectMenuOptionBuilder,
} from "discord.js";
import config from "@/config";
import { randomUUID } from "crypto";

export const data = new SlashCommandSubcommandBuilder()
  .setName("create")
  .setDescription("リアクションパネルを作成します")
  .addRoleOption((option) =>
    option.setName("role0").setDescription("付与するロールを選択してください").setRequired(true)
  )
  .addStringOption((option) =>
    option
      .setName("title")
      .setDescription("役職パネルの名前を設定してください(任意、デフォルトでは役職パネル)")
  )
  .addRoleOption((option) =>
    option.setName("role1").setDescription("付与するロールを選択してください(任意)")
  )
  .addRoleOption((option) =>
    option.setName("role2").setDescription("付与するロールを選択してください(任意)")
  );

export const adminGuildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const member = interaction.member as GuildMember;

  if (!interaction.channel || !interaction.channel.isSendable()) {
    await interaction.reply({
      content: "メッセージを送信できるチャンネルではありません",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }

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

  const panelTitle = interaction.options.getString("title") || "役職パネル";

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

  const panelId = randomUUID();
  const selectMenu = new StringSelectMenuBuilder()
    .setCustomId(`rp_${panelId}`)
    .setPlaceholder("ロールを選択してください")
    .setMinValues(0)
    .setMaxValues(roles.length);

  // 選択肢を追加
  roles.forEach((role) => {
    selectMenu.addOptions(
      new StringSelectMenuOptionBuilder()
        .setLabel(role.name)
        .setValue(role.id.toString())
        .setDescription(`${role.name} ロールを取得/解除します`)
    );
  });
  const row = new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu);

  // パネルの説明を作成
  let description = "下のメニューから希望するロールを選択してください。\n";
  description += "すでに持っているロールを選択すると解除されます。\n\n";

  const send = new EmbedBuilder()
    .setColor("#4CAF50")
    .setTitle(panelTitle)
    .setDescription(description)
    .setTimestamp();

  await interaction.channel.send({
    embeds: [send],
    components: [row],
  });

  const replyEmbed = new EmbedBuilder()
    .setColor("#4CAF50")
    .setTitle("役職パネル作成完了")
    .setDescription("役職パネルが作成されました。")
    .setTimestamp();

  await interaction.reply({
    embeds: [replyEmbed],
    flags: MessageFlags.Ephemeral,
  });
};

export default {
  data,
  adminGuildOnly,
};
