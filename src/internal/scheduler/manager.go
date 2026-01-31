package scheduler

import (
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

type Manager struct {
	ctx  *internal.BotContext
	s    *discordgo.Session
	stop chan struct{}
}

// 新しいスケジューラを作成する
func NewManager(ctx *internal.BotContext, s *discordgo.Session) *Manager {
	return &Manager{ctx: ctx, s: s, stop: make(chan struct{})}
}

// スケジューラを開始する
func (m *Manager) Start() {
	go m.loop()
}

// スケジューラを停止する
func (m *Manager) Stop() {
	close(m.stop)
}

func (m *Manager) loop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.processDue()
		case <-m.stop:
			return
		}
	}
}

func (m *Manager) processDue() {
	now := time.Now().Unix()
	repo := repository.NewScheduleSettingRepository(m.ctx.DB)

	settings, err := repo.ListDue(now)
	if err != nil {
		log.Printf("Failed to load schedules: %v", err)
		return
	}

	for _, setting := range settings {
		m.execute(setting)
	}
}

func (m *Manager) execute(setting *model.ScheduleSetting) {
	if setting == nil {
		return
	}

	_, err := m.s.ChannelMessageSend(setting.ChannelID, setting.Content)
	if err != nil {
		log.Printf("Failed to send scheduled message (id=%s): %v", setting.ID, err)
		return
	}

	repo := repository.NewScheduleSettingRepository(m.ctx.DB)

	// 単発なら削除
	if setting.Cron == "" {
		err = repo.DeleteByID(setting.ID)
		if err != nil {
			log.Printf("Failed to delete schedule (id=%s): %v", setting.ID, err)
		}
		return
	}

	// 繰り返しの場合は次回実行時刻を更新
	base := time.Unix(setting.NextRunAt, 0).In(JST())
	nextRunAt, err := NextRunAtFromCron(setting.Cron, base)
	if err != nil {
		log.Printf("Failed to parse cron (id=%s): %v", setting.ID, err)
		return
	}

	setting.NextRunAt = nextRunAt.Unix()
	err = repo.Update(setting)
	if err != nil {
		log.Printf("Failed to update schedule (id=%s): %v", setting.ID, err)
	}
}

// Cron文字列から次回実行時間を取得する
func NextRunAtFromCron(cronText string, base time.Time) (time.Time, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronText)
	if err != nil {
		return time.Time{}, err
	}

	return schedule.Next(base), nil
}

// JSTロケーション
func JST() *time.Location {
	return time.FixedZone("JST", 9*60*60)
}
