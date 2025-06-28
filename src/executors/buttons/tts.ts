import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonInteraction,
  ButtonStyle,
  EmbedBuilder,
  MessageFlags,
  PermissionsBitField,
  StringSelectMenuBuilder,
  StringSelectMenuInteraction,
} from "discord.js";
import config from "@/config";
import { subSelectMenusHandling } from "@/lib/commandUtils";
import { GetErrorChannel } from "@/lib/channelUtils";
import { uuid58Encode } from "@nakanoaas/uuid58";
import { RPC, AudioLibrary } from "voicevox.js";
import { listTtsDictionary } from "@/lib/dataUtils";

export const name = "tts";

export const execute = async (interaction: ButtonInteraction) => {
  try {
    const id = interaction.customId.split("_");
    const components = [];
    switch (id[1]) {
      case "speaker":
        if (!RPC.rpc) {
          const headers = {
            Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
          };
          await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
        }
        let speakers = await AudioLibrary.getSpeakers();
        if (speakers.length === 0) {
          return interaction.reply({
            content: "No speakers available. Please ask bot admin to add a speaker first.",
            ephemeral: true,
          });
        }
        speakers.sort((a, b) => a.name.localeCompare(b.name));
        if (id.length >= 5 && id[2] == "page") {
          switch (id[3]) {
            case "next":
              if (speakers.length > parseInt(id[4]) * 24) {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + id[4])
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(false),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + (parseInt(id[4]) + 1).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(false),
                ];
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
                speakers = speakers.slice((parseInt(id[4]) - 1) * 24, parseInt(id[4]) * 24);
              } else {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + id[4])
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(false),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + (parseInt(id[4]) + 1).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(true),
                ];
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
                speakers = speakers.slice((parseInt(id[4]) - 1) * 24, speakers.length + 1);
              }
              break;
            case "prev":
              if (parseInt(id[4]) - 1 > 1) {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + (parseInt(id[4]) - 1).toString())
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(false),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + parseInt(id[4]).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(false),
                ];
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
                speakers = speakers.slice((parseInt(id[4]) - 1) * 24, parseInt(id[4]) * 24);
              } else {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + (parseInt(id[4]) - 1).toString())
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(true),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + parseInt(id[4]).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(false),
                ];
                speakers.splice(24);
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
              }
              break;
          }
          const selectMenu = new StringSelectMenuBuilder()
            .setCustomId("tts_speaker_select")
            .setPlaceholder("Select a speaker...")
            .addOptions([
              ...speakers.map((speaker) => ({
                label: speaker.name,
                value: uuid58Encode(speaker.speakerUuid),
              })),
              {
                label: "Cancel",
                value: "cancel",
                description: "Cancel the selection",
              },
            ]);
          selectMenu.setMinValues(1).setMaxValues(1);
          components.push(
            new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu)
          );
          interaction.update({ components });
        }
        break;
      case "dict":
        let allWords = await listTtsDictionary(
          interaction.guild!.id,
          !(interaction.member?.permissions as PermissionsBitField).has(
            PermissionsBitField.Flags.Administrator
          )
            ? interaction.user.id
            : undefined
        );
        if (id.length >= 5 && id[2] == "page") {
          switch (id[3]) {
            case "next":
              if (allWords.length > parseInt(id[4]) * 25) {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + id[4])
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(false),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + (parseInt(id[4]) + 1).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(false),
                ];
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
                allWords = allWords.slice((parseInt(id[4]) - 1) * 24, parseInt(id[4]) * 25);
              } else {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + id[4])
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(false),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + (parseInt(id[4]) + 1).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(true),
                ];
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
                allWords = allWords.slice((parseInt(id[4]) - 1) * 25, allWords.length + 1);
              }
              break;
            case "prev":
              if (parseInt(id[4]) - 1 > 1) {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + (parseInt(id[4]) - 1).toString())
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(false),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + parseInt(id[4]).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(false),
                ];
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
                allWords = allWords.slice((parseInt(id[4]) - 1) * 25, parseInt(id[4]) * 25);
              } else {
                const pagenation = [
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_prev_" + (parseInt(id[4]) - 1).toString())
                    .setLabel("Previous")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("⬅️")
                    .setDisabled(true),
                  new ButtonBuilder()
                    .setCustomId("tts_speaker_page_next_" + parseInt(id[4]).toString())
                    .setLabel("Next")
                    .setStyle(ButtonStyle.Primary)
                    .setEmoji("➡️")
                    .setDisabled(false),
                ];
                allWords.splice(25);
                components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
              }
              break;
          }
          const selectMenu = new StringSelectMenuBuilder()
            .setCustomId("tts_speaker_select")
            .setPlaceholder("Select a speaker...")
            .addOptions([
              ...allWords.map((word) => ({
                label: word.word,
                value: word.id,
              })),
            ]);
          selectMenu.setMinValues(1).setMaxValues(1);
          components.push(
            new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu)
          );
        }
        break;
    }
    interaction.update({ components });
  } catch (error) {
    console.error(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred in ${interaction.customId}\nDetails:\n${error}`
    );
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + (error as any).toString() + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    const channel = await GetErrorChannel(interaction.client);
    if (channel) {
      channel.send({ embeds: [logEmbed] });
    }
    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.error)
      .setTimestamp();
    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({
        embeds: [messageEmbed],
        flags: MessageFlags.SuppressEmbeds,
      });
    }
  }
};
