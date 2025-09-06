import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonInteraction,
  ButtonStyle,
  EmbedBuilder,
  MessageFlags,
  PermissionsBitField,
  StringSelectMenuBuilder,
} from "discord.js";
import config from "@/config";
import { GetErrorChannel } from "@/lib/channelUtils";
import { uuid58Encode } from "@nakanoaas/uuid58";
import { RPC, AudioLibrary } from "voicevox.js";
import { listTtsDictionary } from "@/lib/dataUtils";
import { loggingSystem } from "@/index";

export const name = "tts";

const getPaginationButtons = (
  type: string,
  currentPage: number,
  hasPrev: boolean,
  hasNext: boolean
) => [
  new ButtonBuilder()
    .setCustomId(`tts_${type}_page_prev_${currentPage - 1}`)
    .setLabel("Previous")
    .setStyle(ButtonStyle.Primary)
    .setEmoji("⬅️")
    .setDisabled(!hasPrev),
  new ButtonBuilder()
    .setCustomId(`tts_${type}_page_next_${currentPage + 1}`)
    .setLabel("Next")
    .setStyle(ButtonStyle.Primary)
    .setEmoji("➡️")
    .setDisabled(!hasNext),
];

const getSpeakerSelectMenu = (speakers: any[]) =>
  new StringSelectMenuBuilder()
    .setCustomId("tts_speaker_select")
    .setPlaceholder("話者を選択...")
    .addOptions([
      ...speakers.map((speaker) => ({
        label: speaker.name,
        value: uuid58Encode(speaker.speakerUuid),
      })),
      {
        label: "Cancel",
        value: "cancel",
        description: "選択をキャンセル",
      },
    ])
    .setMinValues(1)
    .setMaxValues(1);

const getDictSelectMenu = (words: any[], action: string) =>
  new StringSelectMenuBuilder()
    .setCustomId(`tts_dict_${action}`)
    .setPlaceholder("単語を選択...")
    .addOptions(
      words.map((word) => ({
        label: word.word,
        value: word.id,
      }))
    )
    .setMinValues(1)
    .setMaxValues(1);

export const execute = async (interaction: ButtonInteraction) => {
  const logger = loggingSystem.getLogger({ function: "TTSButtonInteraction" });
  try {
    const id = interaction.customId.split("_");
    const components: any[] = [];

    if (id[1] === "speaker") {
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

      if (id.length >= 5 && id[2] === "page") {
        const page = parseInt(id[4]);
        const pageSize = 24;
        const total = speakers.length;
        let pagedSpeakers: any[] = [];
        let hasPrev = page > 1;
        let hasNext = total > page * pageSize;

        if (id[3] === "next" || id[3] === "prev") {
          pagedSpeakers = speakers.slice((page - 1) * pageSize, page * pageSize);
          components.push(
            new ActionRowBuilder<ButtonBuilder>().addComponents(
              getPaginationButtons("speaker", page, hasPrev, hasNext)
            )
          );
        }
        components.push(
          new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(
            getSpeakerSelectMenu(pagedSpeakers)
          )
        );
        return interaction.update({ components });
      }
    }

    if (id[1] === "dict") {
      let allWords = await listTtsDictionary(
        interaction.guild!.id,
        !(interaction.member?.permissions as PermissionsBitField).has(
          PermissionsBitField.Flags.Administrator
        )
          ? interaction.user.id
          : undefined
      );

      if (id.length >= 5 && id[3] === "page") {
        const page = parseInt(id[5]);
        const pageSize = 25;
        const total = allWords.length;
        let pagedWords: any[] = [];
        let hasPrev = page > 1;
        let hasNext = total > page * pageSize;

        if (id[4] === "next" || id[4] === "prev") {
          pagedWords = allWords.slice((page - 1) * pageSize, page * pageSize);
          components.push(
            new ActionRowBuilder<ButtonBuilder>().addComponents(
              getPaginationButtons("dict_remove", page, hasPrev, hasNext)
            )
          );
        }
        components.push(
          new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(
            getDictSelectMenu(pagedWords, id[2])
          )
        );
      }
    }

    await interaction.update({ components });
  } catch (error) {
    logger.error(
      { stack_trace: (error as Error).stack },
      "Error in TTS button interaction:",
      error
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
