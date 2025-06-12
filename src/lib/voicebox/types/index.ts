type Mora = {
  text: string;
  consonant: string;
  consonant_length: number;
  vowel: string;
  vowel_length: number;
  pitch: number;
};

type AccentPhrase = {
  moras: Mora[];
  accent: number;
  pause_mora: Mora;
  is_interrogative: boolean;
};

type VoiceSettings = {
  accent_phrases: AccentPhrase[];
  speedScale: number;
  pitchScale: number;
  intonationScale: number;
  volumeScale: number;
  prePhonemeLength: number;
  postPhonemeLength: number;
  pauseLength: number;
  pauseLengthScale: number;
  outputSamplingRate: number;
  outputStereo: boolean;
  kana: string;
};
