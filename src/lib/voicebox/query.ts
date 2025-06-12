import { rpc } from "./rpc";

export const getQuery = async (
  text: string,
  preset?: number,
  speaker?: number,
  enableKatakanaEnglish: boolean = true,
  version?: string
): Promise<VoiceSettings> => {
  if (!rpc) {
    throw new Error("Voicebox API is not connected. Please call connect() first.");
  }
  if (!text || text.trim() === "") {
    throw new Error("Text cannot be empty.");
  }
  if (preset) {
    if (typeof preset !== "number" || preset < 0) {
      throw new Error("Preset must be a non-negative integer.");
    }
    try {
      const response = await rpc.post("/audio_query_from_preset", {
        text,
        preset_id: preset,
        enable_katakana_english: enableKatakanaEnglish,
        core_version: version,
      });
      if (response.status !== 200) {
        throw new Error(
          `Voicebox API returned status code ${response.status}(${response.statusText})`
        );
      }
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch voice settings from Voicebox API: ${error}`);
    }
  } else {
    if (typeof speaker !== "number" || speaker < 0) {
      throw new Error("Speaker must be a non-negative integer.");
    }
    try {
      const response = await rpc.post("/audio_query", {
        text,
        speaker,
        enable_katakana_english: enableKatakanaEnglish,
        core_version: version,
      });
      if (response.status !== 200) {
        throw new Error(
          `Voicebox API returned status code ${response.status}(${response.statusText})`
        );
      }
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch voice settings from Voicebox API: ${error}`);
    }
  }
};
