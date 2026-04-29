package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"
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

func (r *Repository) CheckCurrentDraft(creatorID uint) (ds.SystemLoad, error) {
	if creatorID == 0 {
		return ds.SystemLoad{}, ErrNotAllowed
	}
	var load ds.SystemLoad
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&load)
	if res.Error != nil {
		return ds.SystemLoad{}, res.Error
	}
	if res.RowsAffected == 0 {
		return ds.SystemLoad{}, ErrNoDraft
	}
	return load, nil
}

func (r *Repository) GetSystemLoadDraft(creatorID uint) (ds.SystemLoad, bool, error) {
	load, err := r.CheckCurrentDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
		load = ds.SystemLoad{
			Status:    "draft",
			CreatedAt: time.Now(),
			CreatorID: creatorID,
		}
		if err := r.db.Create(&load).Error; err != nil {
			return ds.SystemLoad{}, false, err
		}
		return load, true, nil
	}
	if err != nil {
		return ds.SystemLoad{}, false, err
	}
	return load, false, nil
}

func (r *Repository) GetModeratorAndCreatorLogin(load ds.SystemLoad) (string, string, error) {
	var creator ds.Users
	if err := r.db.Where("user_id = ?", load.CreatorID).First(&creator).Error; err != nil {
		return "", "", err
	}
	var moderatorLogin string
	if load.ModeratorID != nil && *load.ModeratorID != 0 {
		var moderator ds.Users
		if err := r.db.Where("user_id = ?", *load.ModeratorID).First(&moderator).Error; err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) GetCompletedItemCount(loadID uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.SystemLoadStrategy{}).
		Where("system_load_id = ? AND response_time IS NOT NULL", loadID).
		Count(&count).Error
	return int(count), err
}

func (r *Repository) GetAllSystemLoads(from, to time.Time, status string, creatorID uint) ([]ds.SystemLoad, error) {
	var loads []ds.SystemLoad
	sub := r.db.Where("status != ? AND status != ?", "deleted", "draft")
	if !from.IsZero() {
		sub = sub.Where("forming_date >= ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("forming_date < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	if creatorID != 0 {
		sub = sub.Where("creator_id = ?", creatorID)
	}
	err := sub.Order("system_load_id").Find(&loads).Error
	return loads, err
}

func (r *Repository) GetSingleSystemLoad(id int) (ds.SystemLoad, error) {
	if id < 0 {
		return ds.SystemLoad{}, errors.New("неверное id")
	}
	var load ds.SystemLoad
	err := r.db.Where("system_load_id = ?", id).First(&load).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.SystemLoad{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.SystemLoad{}, err
	}
	if load.Status == "deleted" {
		return ds.SystemLoad{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return load, nil
}

func (r *Repository) GetSystemLoadItems(loadID int) ([]ds.SystemLoadStrategy, error) {
	var items []ds.SystemLoadStrategy
	err := r.db.Where("system_load_id = ?", loadID).
		Preload("Strategy").
		Find(&items).Error
	return items, err
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

	var strategy ds.ShardingStrategy
	if err := r.db.Where("strategy_id = ? AND is_deleted = ?", strategyID, false).First(&strategy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: стратегия с id %d", ErrNotFound, strategyID)
		}
		return err
	}

	var count int64
	r.db.Model(&ds.SystemLoadStrategy{}).
		Where("system_load_id = ? AND strategy_id = ?", load.SystemLoadID, strategyID).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("%w: стратегия %d уже в заявке %d", ErrAlreadyExists, strategyID, load.SystemLoadID)
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
	return r.db.Create(&item).Error
}

func (r *Repository) DeleteStrategyFromSystemLoad(systemLoadID, strategyID int) (ds.SystemLoad, error) {
	var load ds.SystemLoad
	err := r.db.Where("system_load_id = ?", systemLoadID).First(&load).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.SystemLoad{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, systemLoadID)
		}
		return ds.SystemLoad{}, err
	}
	if load.Status != "draft" {
		return ds.SystemLoad{}, fmt.Errorf("%w: можно удалять только из черновика", ErrNotAllowed)
	}
	err = r.db.Where("system_load_id = ? AND strategy_id = ?", systemLoadID, strategyID).
		Delete(&ds.SystemLoadStrategy{}).Error
	if err != nil {
		return ds.SystemLoad{}, err
	}
	return load, nil
}

func (r *Repository) EditStrategyInSystemLoad(systemLoadID, strategyID int, j serializer.SystemLoadStrategyJSON) (ds.SystemLoadStrategy, error) {
	var item ds.SystemLoadStrategy
	err := r.db.Where("system_load_id = ? AND strategy_id = ?", systemLoadID, strategyID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.SystemLoadStrategy{}, fmt.Errorf("%w: стратегия в заявке", ErrNotFound)
		}
		return ds.SystemLoadStrategy{}, err
	}

	var load ds.SystemLoad
	if err := r.db.Where("system_load_id = ?", systemLoadID).First(&load).Error; err != nil {
		return ds.SystemLoadStrategy{}, err
	}
	if load.Status != "draft" {
		return ds.SystemLoadStrategy{}, fmt.Errorf("%w: можно редактировать только черновик", ErrNotAllowed)
	}

	updates := map[string]interface{}{
		"data_volume": j.DataVolume,
		"query_count": j.QueryCount,
	}

	var strategy ds.ShardingStrategy
	if err := r.db.First(&strategy, strategyID).Error; err == nil {
		rt := CalculateResponseTime(
			strategy.LatencyCoefficient,
			strategy.ThroughputCoefficient,
			j.DataVolume, j.QueryCount,
		)
		updates["response_time"] = rt
	}

	err = r.db.Model(&item).Updates(updates).Error
	if err != nil {
		return ds.SystemLoadStrategy{}, err
	}
	r.db.Where("system_load_id = ? AND strategy_id = ?", systemLoadID, strategyID).
		Preload("Strategy").First(&item)
	return item, nil
}

func (r *Repository) EditSystemLoad(id int, j serializer.SystemLoadJSON) (ds.SystemLoad, error) {
	var load ds.SystemLoad
	err := r.db.Where("system_load_id = ? AND status != ?", id, "deleted").First(&load).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.SystemLoad{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.SystemLoad{}, err
	}
	if load.Status != "draft" {
		return ds.SystemLoad{}, fmt.Errorf("%w: можно редактировать только черновик", ErrNotAllowed)
	}
	updates := serializer.SystemLoadFromJSON(j)
	err = r.db.Model(&load).Updates(updates).Error
	if err != nil {
		return ds.SystemLoad{}, err
	}
	r.db.Where("system_load_id = ?", id).First(&load)
	return load, nil
}

func (r *Repository) FormSystemLoad(id int, creatorID uint) (ds.SystemLoad, error) {
	load, err := r.GetSingleSystemLoad(id)
	if err != nil {
		return ds.SystemLoad{}, err
	}
	if load.Status != "draft" {
		return ds.SystemLoad{}, fmt.Errorf("%w: только черновик можно сформировать", ErrNotAllowed)
	}
	if load.CreatorID != creatorID {
		return ds.SystemLoad{}, fmt.Errorf("%w: вы не создатель этой заявки", ErrNotAllowed)
	}

	items, err := r.GetSystemLoadItems(int(load.SystemLoadID))
	if err != nil {
		return ds.SystemLoad{}, err
	}
	if len(items) == 0 {
		return ds.SystemLoad{}, errors.New("нельзя сформировать пустую заявку")
	}

	for _, item := range items {
		var strategy ds.ShardingStrategy
		if err := r.db.First(&strategy, item.StrategyID).Error; err != nil {
			return ds.SystemLoad{}, err
		}
		rt := CalculateResponseTime(
			strategy.LatencyCoefficient,
			strategy.ThroughputCoefficient,
			item.DataVolume, item.QueryCount,
		)
		r.db.Model(&ds.SystemLoadStrategy{}).
			Where("system_load_id = ? AND strategy_id = ?", load.SystemLoadID, item.StrategyID).
			Update("response_time", rt)
	}

	formingDate := time.Now()
	err = r.db.Model(&load).Updates(map[string]interface{}{
		"status":       "formed",
		"forming_date": formingDate,
	}).Error
	if err != nil {
		return ds.SystemLoad{}, err
	}
	load.Status = "formed"
	load.FormingDate = &formingDate
	return load, nil
}

func (r *Repository) FinishSystemLoad(id int, status string, moderatorID uint) (ds.SystemLoad, error) {
	if status != "completed" && status != "rejected" {
		return ds.SystemLoad{}, errors.New("неверный статус: допустимы completed или rejected")
	}
	load, err := r.GetSingleSystemLoad(id)
	if err != nil {
		return ds.SystemLoad{}, err
	}
	if load.Status != "formed" {
		return ds.SystemLoad{}, fmt.Errorf("%w: завершить/отклонить можно только сформированную заявку", ErrNotAllowed)
	}
	finishDate := time.Now()
	err = r.db.Model(&load).Updates(map[string]interface{}{
		"status":       status,
		"finish_date":  finishDate,
		"moderator_id": moderatorID,
	}).Error
	if err != nil {
		return ds.SystemLoad{}, err
	}
	load.Status = status
	load.FinishDate = sql.NullTime{Time: finishDate, Valid: true}
	load.ModeratorID = &moderatorID
	return load, nil
}

func (r *Repository) DeleteSystemLoad(loadID int, creatorID uint) (ds.SystemLoad, error) {
	load, err := r.GetSingleSystemLoad(loadID)
	if err != nil {
		return ds.SystemLoad{}, err
	}
	if load.Status != "draft" {
		return ds.SystemLoad{}, fmt.Errorf("%w: удалить можно только черновик", ErrNotAllowed)
	}
	if load.CreatorID != creatorID {
		return ds.SystemLoad{}, fmt.Errorf("%w: вы не создатель этой заявки", ErrNotAllowed)
	}
	formingDate := time.Now()
	err = r.db.Model(&load).Updates(map[string]interface{}{
		"status":       "deleted",
		"forming_date": formingDate,
	}).Error
	if err != nil {
		return ds.SystemLoad{}, err
	}
	load.Status = "deleted"
	load.FormingDate = &formingDate
	return load, nil
}
