import {
  ActivityOptions,
  ChatInputCommandInteraction,
  CommandInteraction,
  CommandInteractionOptionResolver,
  GuildMemberRoleManager,
  PresenceStatusData,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import { ActivityType } from "discord.js";
import Discord from "discord.js";
import config from "@/config";

export const data = new SlashCommandSubcommandBuilder()
  .setName("setstatus")
  .setDescription("メンテモード")
  .addStringOption((option) =>
    option.setName("statusdiscription").setDescription("プレイ中に表示するやつ(StreamはURLでもok)")
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
      .setChoices({ name: "オン", value: "on" }, { name: "オフ", value: "off" })
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!(interaction.member?.roles as GuildMemberRoleManager).cache.has(config.adminRoleId)) {
    interaction.reply("権限がありません");
    return;
  }
  const client = interaction.client;
  try {
    const onoff = (interaction.options as CommandInteractionOptionResolver).getString("enable");
    if (onoff == "on") {
      const status = (interaction.options as CommandInteractionOptionResolver).getString("status");
      const statusDescription = (interaction.options as CommandInteractionOptionResolver).getString(
        "statusdiscription"
      );
      const activityType = (interaction.options as CommandInteractionOptionResolver).getString(
        "activitytype"
      );
      const activityOpt: ActivityOptions = {
        name: statusDescription as string,
        type: ActivityType.Playing,
      };
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
          //activityOpt.url = statusDescription;
          activityOpt.name = "Youtube";
          break;

        case "CUSTOM":
          activityOpt.type = ActivityType.Custom;
          break;

        default:
          activityOpt.type = ActivityType.Playing;
          break;
      }
      /*
                if (status == 'Discord Android') {
                    client.ws = { properties: { "$os": "Untitled OS", "$browser": "Discord Android", "$device": "Replit Container" } };
                    client.user.setStatus('online');
                } else {*/
      //    client.ws = () => { return { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } }; }
      //}

      client.user.setActivity(activityOpt);
      if (status) {
        client.user.setStatus(status as PresenceStatusData);
      }

      const embed = new Discord.EmbedBuilder()
        .setTitle("ok")
        .setColor(client.config.color.success)
        .setTimestamp();

      interaction.reply({ embeds: [embed] });
      client.function.logUtils.writeConfig(
        {
          onoff: "on",
          status: status,
          playing: statusDescription,
          type: activityType,
        },
        "status"
      );
      return `{ onoff:"on",status: "${status}", statusDesc: "${statusDescription}", type: "${activityType}" }`;
    } else {
      client.user.setActivity(`Servers: ${client.guilds.cache.size}`);
      //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
      client.user.setStatus("online");
      const embed = new Discord.EmbedBuilder()
        .setTitle("ok")
        .setColor(client.config.color.success)
        .setTimestamp();

      interaction.reply({ embeds: [embed] });
      client.function.logUtils.writeConfig({ onoff: "off" }, "status");
      return `{ onoff:"off"}`;
    }
  } catch (e) {
    throw e;
  }
};

export default {
  data,
  execute,
};
