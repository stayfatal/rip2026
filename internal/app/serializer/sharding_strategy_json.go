package serializer

import (
	"strings"

	"web_backend/internal/app/ds"
)

type ShardingStrategyJSON struct {
	StrategyID             uint    `json:"strategy_id"`
	Title                  string  `json:"title"`
	Description            string  `json:"description"`
	ShortDescriptionEN     string  `json:"short_description_en"`
	IsDeleted              bool    `json:"is_deleted"`
	PhotoURL               string  `json:"photo_url"`
	Video                  string  `json:"video"`
	LatencyCoefficient     float64 `json:"latency_coefficient"`
	ThroughputCoefficient  float64 `json:"throughput_coefficient"`
	ReliabilityCoefficient float64 `json:"reliability_coefficient"`
}

func ShardingStrategyToJSON(s ds.ShardingStrategy) ShardingStrategyJSON {
	return ShardingStrategyJSON{
		StrategyID:             s.StrategyID,
		Title:                  s.Title,
		Description:            s.Description,
		ShortDescriptionEN:     s.ShortDescriptionEN,
		IsDeleted:              s.IsDeleted,
		PhotoURL:               s.PhotoURL,
		Video:                  s.Video,
		LatencyCoefficient:     s.LatencyCoefficient,
		ThroughputCoefficient:  s.ThroughputCoefficient,
		ReliabilityCoefficient: s.ReliabilityCoefficient,
	}
}

func ShardingStrategyFromJSON(j ShardingStrategyJSON) ds.ShardingStrategy {
	shortDescriptionEN := strings.TrimSpace(j.ShortDescriptionEN)
	if shortDescriptionEN == "" {
		shortDescriptionEN = "Database sharding strategy profile."
	}
	return ds.ShardingStrategy{
		Title:                  j.Title,
		Description:            j.Description,
		ShortDescriptionEN:     shortDescriptionEN,
		LatencyCoefficient:     j.LatencyCoefficient,
		ThroughputCoefficient:  j.ThroughputCoefficient,
		ReliabilityCoefficient: j.ReliabilityCoefficient,
	}
}
