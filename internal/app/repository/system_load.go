package repository

import (
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CalculateResponseTime(latencyCoeff, throughputCoeff float64, dataVolumeGB, queryCount int) float64 {
	if throughputCoeff <= 0 {
		throughputCoeff = 1
	}
	dataFactor := float64(dataVolumeGB) / 100.0
	queryFactor := float64(queryCount) / 1000.0
	return (dataFactor + queryFactor) * latencyCoeff / throughputCoeff * 10
}

func (r *Repository) GetSystemLoadStrategyCount(creatorID uint) int64 {
	var loadID uint
	err := r.db.Model(&ds.SystemLoad{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("system_load_id").First(&loadID).Error
	if err != nil {
		return 0
	}

	var count int64
	err = r.db.Model(&ds.SystemLoadStrategy{}).
		Where("system_load_id = ?", loadID).Count(&count).Error
	if err != nil {
		logrus.Error("error counting system_load_strategies:", err)
	}
	return count
}

func (r *Repository) GetActiveSystemLoadID(creatorID uint) uint {
	var loadID uint
	err := r.db.Model(&ds.SystemLoad{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Select("system_load_id").First(&loadID).Error
	if err != nil {
		return 0
	}
	return loadID
}

func (r *Repository) GetSystemLoad(id int, creatorID uint) ([]ds.SystemLoadStrategy, *ds.SystemLoad, error) {
	var load ds.SystemLoad
	err := r.db.Where("system_load_id = ? AND creator_id = ? AND status != ?",
		id, creatorID, "deleted").First(&load).Error
	if err != nil {
		return nil, nil, err
	}

	var items []ds.SystemLoadStrategy
	err = r.db.Where("system_load_id = ?", id).
		Preload("Strategy").
		Find(&items).Error
	if err != nil {
		return nil, nil, err
	}
	return items, &load, nil
}

func (r *Repository) AddStrategy(strategyID uint, creatorID uint) error {
	var load ds.SystemLoad

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
		First(&load).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		load = ds.SystemLoad{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&load).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.SystemLoadStrategy{}).
		Where("system_load_id = ? AND strategy_id = ?", load.SystemLoadID, strategyID).
		Count(&count)

	if count == 0 {
		var strategy ds.ShardingStrategy
		if err := r.db.First(&strategy, strategyID).Error; err != nil {
			return err
		}

		rt := CalculateResponseTime(
			strategy.LatencyCoefficient,
			strategy.ThroughputCoefficient,
			100, 1000,
		)

		item := ds.SystemLoadStrategy{
			SystemLoadID: load.SystemLoadID,
			StrategyID:   strategyID,
			DataVolume:   100,
			QueryCount:   1000,
			ResponseTime: &rt,
		}
		if err := r.db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteSystemLoad — логическое удаление заявки через raw SQL UPDATE (без ORM).
func (r *Repository) DeleteSystemLoad(loadID uint) error {
	query := `
		UPDATE system_loads
		SET status = 'deleted'
		WHERE system_load_id = $1;
	`
	result := r.db.Exec(query, loadID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("system_load with id %d not found", loadID)
	}
	return nil
}

func (r *Repository) IsDraftSystemLoad(loadID int, creatorID uint) (bool, error) {
	var load ds.SystemLoad
	err := r.db.Select("status").Where("system_load_id = ? AND creator_id = ?",
		loadID, creatorID).First(&load).Error
	if err != nil {
		return false, err
	}
	return load.Status == "draft", nil
}
