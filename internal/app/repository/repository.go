package repository

import (
	"fmt"
	"strings"
)

// Repository — хранилище данных (Lab 1: данные в массивах, без БД).
type Repository struct {
}

// NewRepository создаёт новый экземпляр репозитория.
func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// ShardingStrategy — услуга: стратегия шардирования с коэффициентами.
type ShardingStrategy struct {
	ID                     int
	Title                  string
	Description            string
	Image                  string  // ключ изображения в Minio
	Video                  string  // ключ видео в Minio (wibes)
	LatencyCoefficient     float64 // коэффициент задержки
	ThroughputCoefficient  float64 // коэффициент пропускной способности
	ReliabilityCoefficient float64 // коэффициент надёжности
}

// SystemLoad — заявка: описание системы для расчёта нагрузки.
type SystemLoad struct {
	ID            int
	Title         string
	Description   string // описание системы текстом
	Strategies    []SystemLoadStrategy
	StrategyCount int
}

// SystemLoadStrategy — связь м-м: данные и запросы, результат — время отклика.
type SystemLoadStrategy struct {
	Strategy     ShardingStrategy
	DataVolume   int     // объём данных (ГБ)
	QueryCount   int     // запросов в секунду
	ResponseTime float64 // результат: время отклика (мс)
}

// GetStrategies возвращает все стратегии шардирования.
func (r *Repository) GetStrategies() ([]ShardingStrategy, error) {
	strategies := []ShardingStrategy{
		{
			ID:                     1,
			Title:                  "Range Sharding",
			Description:            "Разбивка по диапазону ключа. Для временных рядов, логов и быстрых range-запросов.",
			Image:                  "range_sharding.jpg",
			Video:                  "range_sharding.mp4",
			LatencyCoefficient:     1.8,
			ThroughputCoefficient:  0.6,
			ReliabilityCoefficient: 0.70,
		},
		{
			ID:                     2,
			Title:                  "Hash Sharding",
			Description:            "Распределение по хэшу ключа. Равномерная нагрузка, стандарт для OLTP.",
			Image:                  "hash_sharding.jpg",
			Video:                  "hash_sharding.mp4",
			LatencyCoefficient:     1.0,
			ThroughputCoefficient:  1.2,
			ReliabilityCoefficient: 0.90,
		},
		{
			ID:                     3,
			Title:                  "Geo Sharding",
			Description:            "Данные рядом с пользователем по географии. Низкая латентность и локализация (GDPR).",
			Image:                  "geo_sharding.jpg",
			Video:                  "geo_sharding.mp4",
			LatencyCoefficient:     1.4,
			ThroughputCoefficient:  0.8,
			ReliabilityCoefficient: 0.85,
		},
		{
			ID:                     4,
			Title:                  "Directory-Based Sharding",
			Description:            "Каталог «ключ → шард». Максимальная гибкость маршрутизации, мультитенантность.",
			Image:                  "directory_sharding.jpg",
			Video:                  "directory_sharding.mp4",
			LatencyCoefficient:     2.2,
			ThroughputCoefficient:  0.5,
			ReliabilityCoefficient: 0.95,
		},
		{
			ID:                     5,
			Title:                  "Composite Sharding",
			Description:            "Гибрид: гео + хэш (или иные комбинации). Глобальные системы с равномерной нагрузкой.",
			Image:                  "composite_sharding.jpg",
			Video:                  "composite_sharding.mp4",
			LatencyCoefficient:     1.2,
			ThroughputCoefficient:  1.0,
			ReliabilityCoefficient: 0.92,
		},
		{
			ID:                     6,
			Title:                  "Dynamic Sharding",
			Description:            "Авто split/merge шардов по нагрузке. Облачные БД, эластичное масштабирование.",
			Image:                  "dynamic_sharding.jpg",
			Video:                  "dynamic_sharding.mp4",
			LatencyCoefficient:     0.8,
			ThroughputCoefficient:  1.5,
			ReliabilityCoefficient: 0.88,
		},
	}

	if len(strategies) == 0 {
		return nil, fmt.Errorf("массив стратегий пуст")
	}

	return strategies, nil
}

// GetStrategy возвращает стратегию по ID.
func (r *Repository) GetStrategy(id int) (ShardingStrategy, error) {
	strategies, err := r.GetStrategies()
	if err != nil {
		return ShardingStrategy{}, err
	}

	for _, s := range strategies {
		if s.ID == id {
			return s, nil
		}
	}
	return ShardingStrategy{}, fmt.Errorf("стратегия не найдена")
}

// GetStrategiesByTitle возвращает стратегии, содержащие подстроку в названии или описании.
func (r *Repository) GetStrategiesByTitle(query string) ([]ShardingStrategy, error) {
	strategies, err := r.GetStrategies()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var result []ShardingStrategy
	for _, s := range strategies {
		if strings.Contains(strings.ToLower(s.Title), q) {
			result = append(result, s)
		}
	}
	return result, nil
}

// CalculateResponseTime вычисляет время отклика по коэффициентам стратегии и параметрам нагрузки.
func CalculateResponseTime(latencyCoeff, throughputCoeff float64, dataVolumeGB, queryCount int) float64 {
	if throughputCoeff <= 0 {
		throughputCoeff = 1
	}
	dataFactor := float64(dataVolumeGB) / 100.0
	queryFactor := float64(queryCount) / 1000.0
	return (dataFactor + queryFactor) * latencyCoeff / throughputCoeff * 10
}

// buildSystemLoad собирает заявку из записей м-м.
func (r *Repository) buildSystemLoad(id int, title, description string, entries []struct {
	StrategyID int
	DataVolume int
	QueryCount int
}) (SystemLoad, error) {
	strategies, err := r.GetStrategies()
	if err != nil {
		return SystemLoad{}, err
	}

	stratMap := make(map[int]ShardingStrategy)
	for _, s := range strategies {
		stratMap[s.ID] = s
	}

	var loadStrategies []SystemLoadStrategy

	for _, e := range entries {
		strat, ok := stratMap[e.StrategyID]
		if !ok {
			continue
		}
		rt := CalculateResponseTime(strat.LatencyCoefficient, strat.ThroughputCoefficient, e.DataVolume, e.QueryCount)
		loadStrategies = append(loadStrategies, SystemLoadStrategy{
			Strategy:     strat,
			DataVolume:   e.DataVolume,
			QueryCount:   e.QueryCount,
			ResponseTime: rt,
		})
	}

	return SystemLoad{
		ID:            id,
		Title:         title,
		Description:   description,
		Strategies:    loadStrategies,
		StrategyCount: len(loadStrategies),
	}, nil
}

// GetSystemLoads возвращает все заявки на расчёт нагрузки.
func (r *Repository) GetSystemLoads() ([]SystemLoad, error) {
	entries := []struct {
		StrategyID int
		DataVolume int
		QueryCount int
	}{
		{1, 200, 3000},
		{2, 500, 8000},
		{3, 150, 2000},
		{4, 100, 1500},
		{5, 300, 5000},
		{6, 800, 12000},
	}

	load, err := r.buildSystemLoad(
		1,
		"Нагрузка e-commerce платформы",
		"Интернет-магазин с 50 млн активных пользователей, обрабатывающий каталог товаров, заказы и пользовательские сессии. Система работает в трёх регионах (EU, US, Asia) с требованиями к низкой латентности и высокой доступности. Пиковые нагрузки приходятся на сезонные распродажи.",
		entries,
	)
	if err != nil {
		return nil, err
	}

	return []SystemLoad{load}, nil
}

// GetSystemLoad возвращает заявку по ID.
func (r *Repository) GetSystemLoad(id int) (SystemLoad, error) {
	loads, err := r.GetSystemLoads()
	if err != nil {
		return SystemLoad{}, err
	}

	for _, l := range loads {
		if l.ID == id {
			return l, nil
		}
	}
	return SystemLoad{}, fmt.Errorf("заявка не найдена")
}

// GetSystemLoadForStrategy ищет заявку, содержащую данную стратегию, и возвращает запись м-м.
func (r *Repository) GetSystemLoadForStrategy(strategyID int) (*SystemLoadStrategy, error) {
	loads, err := r.GetSystemLoads()
	if err != nil {
		return nil, err
	}

	for _, l := range loads {
		for _, ls := range l.Strategies {
			if ls.Strategy.ID == strategyID {
				return &ls, nil
			}
		}
	}
	return nil, fmt.Errorf("стратегия не найдена в заявках")
}
