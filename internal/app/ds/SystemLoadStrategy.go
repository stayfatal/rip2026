package ds

type SystemLoadStrategy struct {
	SystemLoadID uint `gorm:"primaryKey;column:system_load_id"`
	StrategyID   uint `gorm:"primaryKey;column:strategy_id"`

	DataVolume   int      `gorm:"not null;default:0"`
	QueryCount   int      `gorm:"not null;default:0"`
	ResponseTime *float64 `gorm:"type:numeric(12,2)"`

	Strategy   ShardingStrategy `gorm:"foreignKey:StrategyID"`
	SystemLoad SystemLoad       `gorm:"foreignKey:SystemLoadID"`
}

func (SystemLoadStrategy) TableName() string {
	return "system_load_strategies"
}
