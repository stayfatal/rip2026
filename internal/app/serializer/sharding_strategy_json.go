package serializer

import "web_backend/internal/app/ds"

type ShardingStrategyJSON struct {
	StrategyID             uint    `json:"strategy_id"`
	Title                  string  `json:"title"`
	Description            string  `json:"description"`
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
		IsDeleted:              s.IsDeleted,
		PhotoURL:               s.PhotoURL,
		Video:                  s.Video,
		LatencyCoefficient:     s.LatencyCoefficient,
		ThroughputCoefficient:  s.ThroughputCoefficient,
		ReliabilityCoefficient: s.ReliabilityCoefficient,
	}
}

func ShardingStrategyFromJSON(j ShardingStrategyJSON) ds.ShardingStrategy {
	return ds.ShardingStrategy{
		Title:                  j.Title,
		Description:            j.Description,
		LatencyCoefficient:     j.LatencyCoefficient,
		ThroughputCoefficient:  j.ThroughputCoefficient,
		ReliabilityCoefficient: j.ReliabilityCoefficient,
	}
}
