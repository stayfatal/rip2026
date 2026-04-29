package serializer

import "web_backend/internal/app/ds"

type SystemLoadStrategyJSON struct {
	SystemLoadID uint     `json:"system_load_id"`
	StrategyID   uint     `json:"strategy_id"`
	DataVolume   int      `json:"data_volume"`
	QueryCount   int      `json:"query_count"`
	ResponseTime *float64 `json:"response_time"`
}

type SystemLoadStrategyDetailJSON struct {
	SystemLoadID uint                 `json:"system_load_id"`
	StrategyID   uint                 `json:"strategy_id"`
	DataVolume   int                  `json:"data_volume"`
	QueryCount   int                  `json:"query_count"`
	ResponseTime *float64             `json:"response_time"`
	Strategy     ShardingStrategyJSON `json:"strategy"`
}

func SystemLoadStrategyToJSON(item ds.SystemLoadStrategy) SystemLoadStrategyJSON {
	return SystemLoadStrategyJSON{
		SystemLoadID: item.SystemLoadID,
		StrategyID:   item.StrategyID,
		DataVolume:   item.DataVolume,
		QueryCount:   item.QueryCount,
		ResponseTime: item.ResponseTime,
	}
}

func SystemLoadStrategyDetailToJSON(item ds.SystemLoadStrategy) SystemLoadStrategyDetailJSON {
	return SystemLoadStrategyDetailJSON{
		SystemLoadID: item.SystemLoadID,
		StrategyID:   item.StrategyID,
		DataVolume:   item.DataVolume,
		QueryCount:   item.QueryCount,
		ResponseTime: item.ResponseTime,
		Strategy:     ShardingStrategyToJSON(item.Strategy),
	}
}

func SystemLoadStrategyFromJSON(j SystemLoadStrategyJSON) ds.SystemLoadStrategy {
	return ds.SystemLoadStrategy{
		DataVolume: j.DataVolume,
		QueryCount: j.QueryCount,
	}
}
