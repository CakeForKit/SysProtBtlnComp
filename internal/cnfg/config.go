package cnfg

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type InputRaceConfig struct {
	Laps        int    `mapstructure:"laps"`
	LapLen      int    `mapstructure:"lapLen"`
	PenaltyLen  int    `mapstructure:"penaltyLen"`
	FiringLines int    `mapstructure:"firingLines"`
	Start       string `mapstructure:"start"`
	StartDelta  string `mapstructure:"startDelta"`
}

type RaceConfig struct {
	Laps        int
	LapLen      int
	PenaltyLen  int
	FiringLines int
	Start       string
	StartDelta  time.Duration
}

var (
	ErrConfigRead        = errors.New("ReadInConfig")
	ErrInvalidTimeFormat = errors.New("invalid time format")
)

func LoadRaceConfig(path string) (*RaceConfig, error) {
	var rawConfig InputRaceConfig
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err := v.Unmarshal(&rawConfig); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	duration, err := parseHHMMSS(rawConfig.StartDelta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse startDelta: %w", err)
	}

	return &RaceConfig{
		Laps:        rawConfig.Laps,
		LapLen:      rawConfig.LapLen,
		PenaltyLen:  rawConfig.PenaltyLen,
		FiringLines: rawConfig.FiringLines,
		Start:       rawConfig.Start,
		StartDelta:  duration,
	}, nil
}

func parseHHMMSS(timeStr string) (time.Duration, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, ErrInvalidTimeFormat
	}

	var h, m, s time.Duration
	var err error

	if h, err = time.ParseDuration(parts[0] + "h"); err != nil {
		return 0, err
	}
	if m, err = time.ParseDuration(parts[1] + "m"); err != nil {
		return 0, err
	}
	if s, err = time.ParseDuration(parts[2] + "s"); err != nil {
		return 0, err
	}

	return h + m + s, nil
}
