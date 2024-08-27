const { SlashCommandBuilder } = require("discord.js");
const { ActivityType } = require("discord-api-types/v9");
const Discord = require("discord.js");
module.exports = {
  guildOnly: false,
  adminGuildOnly: true,
  data: new SlashCommandBuilder()
    .setName("setstatus")
    .setDescription("メンテモード")
    .addStringOption((option) =>
      option
        .setName("statusdiscription")
        .setDescription("プレイ中に表示するやつ(StreamはURLでもok)")
    )
    .addStringOption((option) =>
      option.setName("status").setDescription("すてーたす").setChoices(
        { name: "オンライン", value: "online" },
        { name: "取り込み中", value: "dnd" },
        { name: "スリープ", value: "idle" },
        //{ name: 'スマホでオンライン', value: 'Discord Android' },
        { name: "オンライン隠し", value: "invisible" }
      )
    )
    .addStringOption((option) =>
      option
        .setName("activitytype")
        .setDescription("あくてぃびてぃ")
        .setChoices(
          { name: "視聴中", value: "WATCHING" },
          { name: "プレイ中", value: "PLAYING" },
          { name: "参戦中", value: "COMPETING" },
          { name: "再生中(聞く)", value: "LISTENING" },
          { name: "配信中", value: "STREAMING" },
          { name: "カスタム", value: "CUSTOM" }
        )
    )
    .addStringOption((option) =>
      option
        .setName("enable")
        .setDescription("オンオフ")
        .setChoices(
          { name: "オン", value: "on" },
          { name: "オフ", value: "off" }
        )
    ),

  /**
   * Executes the maintenance command.
   *
   * @param {Interaction} i - The interaction object.
   * @param {Client} client - The Discord client object.
   * @returns {string} - The result of the execution.
   * @throws {Error} - If an error occurs during execution.
   */
  async execute(i) {
    const client = i.client;
    try {
      const onoff = i.options.getString("enable");
      if (onoff == "on") {
        const status = i.options.getString("status");
        const statusDescription = i.options.getString("statusdiscription");
        const activityType = i.options.getString("activitytype");
        let activityOpt = { type: null, name: statusDescription };
        switch (activityType) {
          case "WATCHING":
            activityOpt.type = ActivityType.Watching;
            break;

          case "COMPETING":
            activityOpt.type = ActivityType.Competing;
            break;

          case "LISTENING":
            activityOpt.type = ActivityType.Listening;
            break;

          case "STREAMING":
            activityOpt.type = ActivityType.Streaming;
            activityOpt.url = statusDescription;
            activityOpt.name = "Youtube";
            break;

          case "CUSTOM":
            activityOpt.type = ActivityType.Custom;
            console.log("custom");
            break;

          default:
            activityOpt.type = ActivityType.Playing;
            break;
        }
        console.log(activityOpt.type);
        /*
                if (status == 'Discord Android') {
                    client.ws = { properties: { "$os": "Untitled OS", "$browser": "Discord Android", "$device": "Replit Container" } };
                    client.user.setStatus('online');
                } else {*/
        //    client.ws = () => { return { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } }; }
        //}

        client.user.setActivity(activityOpt);
        client.user.setStatus(status);

        const embed = new Discord.EmbedBuilder()
          .setTitle("ok")
          .setColor(client.conf.color.s)
          .setTimestamp();

        i.reply({ embeds: [embed] });
        client.func.logUtils.loging(
          {
            onoff: "on",
            status: status,
            playing: statusDescription,
            type: activityType,
          },
          "v1/conf/status"
        );
        return `{ onoff:"on",status: "${status}", statusDesc: "${statusDescription}", type: "${activityType}" }`;
      } else {
        client.user.setActivity(`Servers: ${client.guilds.cache.size}`);
        //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
        client.user.setStatus("online");
        const embed = new Discord.EmbedBuilder()
          .setTitle("ok")
          .setColor(client.conf.color.s)
          .setTimestamp();

        i.reply({ embeds: [embed] });
        client.func.logUtils.loging({ onoff: "off" }, "v1/conf/status");
        return `{ onoff:"off"}`;
      }
    } catch (e) {
      throw e;
    }
  },
};
