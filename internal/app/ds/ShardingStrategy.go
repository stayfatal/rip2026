package ds

type ShardingStrategy struct {
	StrategyID             uint    `gorm:"primaryKey;column:strategy_id"`
	Title                  string  `gorm:"type:varchar(255);not null"`
	Description            string  `gorm:"type:varchar(1000);not null"`
	ShortDescriptionEN     string  `gorm:"column:short_description_en;type:varchar(100);not null;default:''"`
	IsDeleted              bool    `gorm:"type:boolean;not null;default:false"`
	PhotoURL               string  `gorm:"column:photo_url;type:varchar(255)"`
	Video                  string  `gorm:"type:varchar(255)"`
	LatencyCoefficient     float64 `gorm:"type:numeric(5,2);not null;default:1.0"`
	ThroughputCoefficient  float64 `gorm:"type:numeric(5,2);not null;default:1.0"`
	ReliabilityCoefficient float64 `gorm:"type:numeric(5,2);not null;default:0.9"`
}

func (ShardingStrategy) TableName() string {
	return "sharding_strategies"
}
