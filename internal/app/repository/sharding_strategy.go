package repository

import (
	"fmt"
	"web_backend/internal/app/ds"
)

func (r *Repository) GetStrategies() ([]ds.ShardingStrategy, error) {
	var strategies []ds.ShardingStrategy
	err := r.db.Where("is_deleted = ?", false).Find(&strategies).Error
	if err != nil {
		return nil, err
	}
	if len(strategies) == 0 {
		return nil, fmt.Errorf("массив стратегий пуст")
	}
	return strategies, nil
}

func (r *Repository) GetStrategy(id int) (ds.ShardingStrategy, error) {
	var strategy ds.ShardingStrategy
	err := r.db.Where("strategy_id = ? AND is_deleted = ?", id, false).First(&strategy).Error
	if err != nil {
		return ds.ShardingStrategy{}, err
	}
	return strategy, nil
}

func (r *Repository) GetStrategiesByTitle(title string) ([]ds.ShardingStrategy, error) {
	var strategies []ds.ShardingStrategy
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&strategies).Error
	if err != nil {
		return nil, err
	}
	return strategies, nil
}
