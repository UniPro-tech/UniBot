import {
  ChatInputCommandInteraction,
  SlashCommandSubcommandBuilder,
  EmbedBuilder,
  GuildMember,
  MessageFlags,
} from "discord.js";
import config from "@/config";

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
      .setRequired(false)
  )
  .addRoleOption((option) =>
    option
      .setName("role1")
      .setDescription("付与するロールを選択してください(任意)")
      .setRequired(false)
  )
  .addRoleOption((option) =>
    option
      .setName("role2")
      .setDescription("付与するロールを選択してください(任意)")
      .setRequired(false)
  );

export const adminGuildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const member = interaction.member as GuildMember;
  if (!member.roles.cache.has(config.adminRoleId)) {
    await interaction.reply({ content: "権限がありません", ephemeral: true });
    return; // アドミンロールが付与されていなかったら終了
  }

  let result = "";
  let roleID_list: string[] = [];
  const alphabet = [
    "🇦",
    "🇧",
    "🇨",
    "🇩",
    "🇪",
    "🇫",
    "🇬",
    "🇭",
    "🇮",
    "🇯",
    "🇰",
    "🇱",
    "🇲",
    "🇳",
    "🇴",
    "🇵",
    "🇶",
    "🇷",
    "🇸",
    "🇹",
    "🇺",
    "🇻",
    "🇼",
    "🇽",
    "🇾",
    "🇿",
  ];

  const memberRoles = member.roles.cache.map((role) => role.position);
  const highestMemberRole = Math.max(...memberRoles);

  const botMember = interaction.guild?.members.me;
  if (!botMember) {
    await interaction.reply({ content: "ボットのメンバー情報が取得できません", ephemeral: true });
    return;
  }
  const botRoles = botMember.roles.cache.map((role) => role.position);
  const highestBotRole = Math.max(...botRoles);

  const panelTitle = interaction.options.getString("title") || "役職パネル";

  for (let i = 0; i <= 10; i++) {
    const role = interaction.options.getRole(`role${i}`);
    if (role) {
      // @everyone ロールのIDを取得
      if (role.id === interaction.guild?.id) {
        await interaction.reply({
          content: "`@everyone` ロールは選択できません。",
          ephemeral: true,
        });
        return;
      }

      // ユーザーの役職よりも高い権限のロールを指定した場合
      if (role.position > highestMemberRole) {
        await interaction.reply({
          content: `指定されたロール ${role.name} はあなたより高い権限を持っています。これを付与することはできません。`,
          ephemeral: true,
        });
        return;
      }

      // ボットの役職よりも高い権限のロールを指定した場合
      if (role.position >= highestBotRole) {
        await interaction.reply({
          content: `指定されたロール ${role.name} はこのボットより高い権限を持っています。これを付与することはできません。`,
          ephemeral: true,
        });
        return;
      }

      roleID_list.push(role.id);
      result += `${alphabet[i]}:<@&${role.id}>\n`;
    }
  }

  if (!interaction.channel || !("send" in interaction.channel)) {
    await interaction.reply({
      content: "メッセージを送信できるチャンネルではありません",
      ephemeral: true,
    });
    return;
  }

  const send = new EmbedBuilder()
    .setColor("#4CAF50")
    .setTitle(panelTitle)
    .setDescription(result)
    .setTimestamp();

  const message = await interaction.channel.send({
    embeds: [send],
  });

  const reply = new EmbedBuilder()
    .setColor("#4CAF50")
    .setTitle("役職パネル作成完了")
    .setDescription("役職パネルが作成されました。")
    .setTimestamp();

  await interaction.reply({
    embeds: [reply],
    flags: MessageFlags.Ephemeral,
  });

  for (let i = 0; i < roleID_list.length; i++) {
    await message.react(alphabet[i]);
  }
};

export default {
  data,
  adminGuildOnly,
};
