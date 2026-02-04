package voicevox

// Speaker は /speakers のレスポンスを表します
type Speaker struct {
	Name        string         `json:"name"`
	Styles      []SpeakerStyle `json:"styles"`
	SpeakerUUID string         `json:"speaker_uuid"`
	Version     string         `json:"version"`
}

// SpeakerStyle は話者のスタイル情報です
type SpeakerStyle struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
