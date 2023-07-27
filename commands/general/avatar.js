const { SlashCommandBuilder } = require("discord.js");
const Discord = require("discord.js");

module.exports = {
  data: new SlashCommandBuilder()
    .setName("avatar")
    .setDescription(
      "アイコンのURLを取得します。"
    )
    .addUserOption((option) =>
      option.setName("target").setDescription("ここにユーザーを指定してそのユーザーのアイコンのURLを取得します。")
    ),
  async execute(i) {
    const user = i.options.getUser("target");
    if (user){
      const embed = new Discord.EmbedBuilder()
      .setTitle(`${user.username}'s Avatar`)
      .setDescription(`URL:${user.displayAvatarURL({ dynamic: true })}`)
      .setImage('https://i.imgur.com/AfFp7pu.png')
      .setColor(client.conf.color.s)
      .setTimestamp();
      i.reply({ embeds: [embed] });
    }
    else{
      const embed = new Discord.EmbedBuilder()
      .setTitle(`${user.username}'s Avatar`)
      .setDescription(`URL:${user.displayAvatarURL({ dynamic: true })}`)
      .setImage(`${user.displayAvatarURL({ dynamic: true })}`)
      .setColor(client.conf.color.s)
      .setTimestamp();
      i.reply({ embeds: [embed] });
    }
    return 'No data';
  },
};
