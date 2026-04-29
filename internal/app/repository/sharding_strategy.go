package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"web_backend/internal/app/ds"
	minioClient "web_backend/internal/app/minioClient"
	"web_backend/internal/app/serializer"
)

func (r *Repository) GetStrategies() ([]ds.ShardingStrategy, error) {
	var strategies []ds.ShardingStrategy
	err := r.db.Where("is_deleted = ?", false).Find(&strategies).Error
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

func (r *Repository) GetStrategy(id int) (*ds.ShardingStrategy, error) {
	var strategy ds.ShardingStrategy
	err := r.db.Where("strategy_id = ? AND is_deleted = ?", id, false).First(&strategy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: стратегия с id %d", ErrNotFound, id)
		}
		return nil, err
	}
	return &strategy, nil
}

func (r *Repository) GetStrategiesByTitle(title string) ([]ds.ShardingStrategy, error) {
	var strategies []ds.ShardingStrategy
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+title+"%", false).Find(&strategies).Error
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

func (r *Repository) CreateStrategy(j serializer.ShardingStrategyJSON) (ds.ShardingStrategy, error) {
	strategy := serializer.ShardingStrategyFromJSON(j)
	err := r.db.Create(&strategy).Scan(&strategy).Error
	if err != nil {
		return ds.ShardingStrategy{}, err
	}
	return strategy, nil
}

func (r *Repository) AddPhoto(ctx *gin.Context, strategyID int, file *multipart.FileHeader) (*ds.ShardingStrategy, error) {
	strategy, err := r.GetStrategy(strategyID)
	if err != nil {
		return nil, err
	}
	if strategy.PhotoURL != "" {
		_ = minioClient.DeleteObject(ctx, r.mc, minioClient.GetImgBucket(), strategy.PhotoURL)
	}
	fileName, err := minioClient.UploadImage(ctx, r.mc, minioClient.GetImgBucket(), file, strategy.StrategyID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&ds.ShardingStrategy{}).Where("strategy_id = ?", strategyID).Update("photo_url", fileName).Error; err != nil {
		return nil, err
	}
	strategy.PhotoURL = fileName
	return strategy, nil
}

func (r *Repository) AddVideo(ctx *gin.Context, strategyID int, file *multipart.FileHeader) (*ds.ShardingStrategy, error) {
	strategy, err := r.GetStrategy(strategyID)
	if err != nil {
		return nil, err
	}
	if strategy.Video != "" {
		_ = minioClient.DeleteObject(ctx, r.mc, minioClient.GetImgBucket(), strategy.Video)
	}
	fileName, err := minioClient.UploadVideo(ctx, r.mc, minioClient.GetImgBucket(), file, strategy.StrategyID)
	if err != nil {
		return nil, err
	}
	if err := r.db.Model(&ds.ShardingStrategy{}).Where("strategy_id = ?", strategyID).Update("video", fileName).Error; err != nil {
		return nil, err
	}
	strategy.Video = fileName
	return strategy, nil
}

func (r *Repository) DeleteStrategy(id int) error {
	strategy, err := r.GetStrategy(id)
	if err != nil {
		return err
	}
	if strategy.PhotoURL != "" {
		_ = minioClient.DeleteObject(context.Background(), r.mc, minioClient.GetImgBucket(), strategy.PhotoURL)
	}
	if strategy.Video != "" {
		_ = minioClient.DeleteObject(context.Background(), r.mc, minioClient.GetImgBucket(), strategy.Video)
	}
	return r.db.Model(&ds.ShardingStrategy{}).Where("strategy_id = ?", id).Update("is_deleted", true).Error
}
