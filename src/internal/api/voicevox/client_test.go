package voicevox_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"unibot/internal/api/voicevox"
	"unibot/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	println("Loading .env file...")
	util.LoadProjectEnv()

	code := m.Run()

	os.Exit(code)
}

// clientが正常に作成されるか
func TestNew(t *testing.T) {
	uri := os.Getenv("VOICEVOX_URI")
	apiKey := os.Getenv("VOICEVOX_API_KEY")
	client := voicevox.New(uri, apiKey)
	assert.Equal(t, uri, client.BaseURL)
	assert.Equal(t, apiKey, client.APIKey)
	assert.NotNil(t, client.HTTP)
}

func TestSynthesize(t *testing.T) {
	client := createNewVoicevoxClient()
	cCtx := context.Background()

	bytes, err := client.Synthesize(cCtx, "このメッセージはテストです。", "0", float64(100)/100.0)
	assert.NoError(t, err)

	tmpPath := filepath.Join(os.TempDir(), "test.wav")

	err = os.WriteFile(tmpPath, bytes, 0644)
	assert.NoError(t, err)

	t.Logf("saved to %s, pls check the file", tmpPath)
}

func createNewVoicevoxClient() *voicevox.Client {
	uri := os.Getenv("VOICEVOX_URI")
	apiKey := os.Getenv("VOICEVOX_API_KEY")
	client := voicevox.New(uri, apiKey)
	return client
}
