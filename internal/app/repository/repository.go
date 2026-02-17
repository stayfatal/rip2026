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

// ShardingStrategy — услуга: стратегия шардирования данных в распределённой системе.
type ShardingStrategy struct {
	ID              int
	Title           string
	Description     string
	Image           string // ключ изображения в Minio
	KeyType         string // тип ключа шардирования
	Uniformity      string // равномерность распределения
	UseCase         string // область применения
	AvgResponseTime int    // базовое время отклика (мс)
}

// LoadCalculation — заявка: расчёт нагрузки для набора стратегий.
type LoadCalculation struct {
	ID                     int
	Title                  string
	Description            string
	Status                 string
	DataVolume             int     // объём данных (ГБ)
	QueryCount             int     // запросов в секунду
	ResultResponseTime     float64 // итоговое время отклика (мс)
	ResultLoadDistribution string  // распределение нагрузки
	Strategies             []CalculationStrategy
	StrategyCount          int
}

// CalculationStrategy — связь м-м: стратегия внутри заявки.
type CalculationStrategy struct {
	Strategy   ShardingStrategy
	ShardCount int    // количество шардов (поле м-м)
	Priority   int    // приоритет
	Comment    string // комментарий пользователя
}

// GetStrategies возвращает все стратегии шардирования.
func (r *Repository) GetStrategies() ([]ShardingStrategy, error) {
	strategies := []ShardingStrategy{
		{
			ID:              1,
			Title:           "Range Sharding",
			Description:     "Range Sharding (шардирование по диапазону) — стратегия распределения данных, при которой записи разбиваются на непрерывные интервалы по значению ключа шардирования. Каждому шарду назначается определённый диапазон ключей, например: записи с ID от 1 до 1 000 000 хранятся на первом шарде, от 1 000 001 до 2 000 000 — на втором и так далее. Это позволяет эффективно выполнять range-запросы, поскольку данные расположены последовательно. Однако при неравномерном распределении записей могут возникнуть хотспоты — ситуации, когда один шард получает непропорционально большую нагрузку. Range Sharding хорошо подходит для временных рядов, логов и данных с естественным порядком.",
			Image:           "range_sharding.jpg",
			KeyType:         "Диапазон",
			Uniformity:      "Низкая",
			UseCase:         "Временные ряды, логи, аналитика",
			AvgResponseTime: 45,
		},
		{
			ID:              2,
			Title:           "Hash Sharding",
			Description:     "Hash Sharding (хэш-шардирование) — стратегия, при которой данные распределяются по шардам на основе хэш-функции, применённой к ключу шардирования. Хэш-значение определяет, на какой шард попадёт запись. Это обеспечивает равномерное распределение данных и нагрузки, что устраняет проблему хотспотов. Однако при использовании хэш-шардирования теряется возможность эффективно выполнять range-запросы, так как соседние по ключу записи могут оказаться на разных шардах. Hash Sharding — это стандартный выбор для систем общего назначения, где важна равномерная загрузка узлов.",
			Image:           "hash_sharding.jpg",
			KeyType:         "Хэш",
			Uniformity:      "Высокая",
			UseCase:         "Системы общего назначения, OLTP",
			AvgResponseTime: 25,
		},
		{
			ID:              3,
			Title:           "Geo Sharding",
			Description:     "Geo Sharding (географическое шардирование) — стратегия распределения данных по географическому признаку. Записи направляются на шарды, расположенные ближе к конечному пользователю, что минимизирует задержку при обращении к данным. Например, данные европейских пользователей хранятся на серверах в Европе, а азиатских — в Азии. Эта стратегия критически важна для глобальных приложений с требованиями к низкой латентности и соблюдению законодательства о локализации данных (GDPR). Главный недостаток — сложность перебалансировки при изменении географии пользовательской базы.",
			Image:           "geo_sharding.jpg",
			KeyType:         "Географический",
			Uniformity:      "Средняя",
			UseCase:         "Геоданные, CDN, глобальные сервисы",
			AvgResponseTime: 35,
		},
		{
			ID:              4,
			Title:           "Directory-Based Sharding",
			Description:     "Directory-Based Sharding (шардирование на основе каталога) — стратегия, при которой используется отдельная таблица-каталог для хранения соответствий между ключом шардирования и конкретным шардом. При каждом запросе система обращается к каталогу, чтобы определить местоположение данных. Это даёт максимальную гибкость: записи можно перемещать между шардами без изменения логики приложения. Однако каталог становится единой точкой отказа (SPOF) и потенциальным узким местом производительности. Directory-Based Sharding применяется в системах со сложными паттернами доступа, где другие стратегии не подходят.",
			Image:           "directory_sharding.jpg",
			KeyType:         "Каталог",
			Uniformity:      "Высокая",
			UseCase:         "Мультитенантные системы, гибкое распределение",
			AvgResponseTime: 55,
		},
		{
			ID:              5,
			Title:           "Composite Sharding",
			Description:     "Composite Sharding (составное шардирование) — гибридная стратегия, комбинирующая два или более подхода к шардированию. Например, данные сначала разделяются по географии (Geo Sharding), а затем внутри каждого региона — по хэшу (Hash Sharding). Это позволяет получить преимущества обоих подходов: низкую латентность за счёт локальности данных и равномерную нагрузку внутри региона. Composite Sharding требует более сложной инфраструктуры и конфигурации, но обеспечивает наилучшую производительность для крупномасштабных глобальных систем с разнородными паттернами доступа.",
			Image:           "composite_sharding.jpg",
			KeyType:         "Составной",
			Uniformity:      "Высокая",
			UseCase:         "Крупные глобальные платформы, микросервисы",
			AvgResponseTime: 30,
		},
		{
			ID:              6,
			Title:           "Dynamic Sharding",
			Description:     "Dynamic Sharding (динамическое шардирование) — стратегия, при которой количество и границы шардов автоматически адаптируются к текущей нагрузке и объёму данных. Система мониторит размер шардов и при превышении порога автоматически разделяет перегруженный шард на два (split), а при снижении нагрузки — объединяет шарды (merge). Это обеспечивает автоматическое масштабирование без ручного вмешательства. Dynamic Sharding используется в облачных базах данных, таких как MongoDB Atlas и Google Cloud Spanner, где эластичность является ключевым требованием.",
			Image:           "dynamic_sharding.jpg",
			KeyType:         "Автоматический",
			Uniformity:      "Высокая",
			UseCase:         "Облачные БД, эластичное масштабирование",
			AvgResponseTime: 20,
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

// GetStrategiesByTitle возвращает стратегии, содержащие подстроку в любом текстовом поле.
func (r *Repository) GetStrategiesByTitle(query string) ([]ShardingStrategy, error) {
	strategies, err := r.GetStrategies()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var result []ShardingStrategy
	for _, s := range strategies {
		if strings.Contains(strings.ToLower(s.Title), q) ||
			strings.Contains(strings.ToLower(s.Description), q) ||
			strings.Contains(strings.ToLower(s.KeyType), q) ||
			strings.Contains(strings.ToLower(s.UseCase), q) ||
			strings.Contains(strings.ToLower(s.Uniformity), q) {
			result = append(result, s)
		}
	}
	return result, nil
}

// CalculateResponseTime вычисляет итоговое время отклика на основе параметров стратегии.
func CalculateResponseTime(baseTime int, shardCount int, dataVolumeGB int) float64 {
	if shardCount <= 0 {
		shardCount = 1
	}
	dataFactor := float64(dataVolumeGB) / 100.0
	return float64(baseTime) * (1 + dataFactor) / float64(shardCount)
}

// buildCalculation собирает заявку из записей.
func (r *Repository) buildCalculation(id int, title string, description string, status string, dataVolume int, queryCount int, entries []struct {
	StrategyID int
	ShardCount int
	Priority   int
	Comment    string
}) (LoadCalculation, error) {
	strategies, err := r.GetStrategies()
	if err != nil {
		return LoadCalculation{}, err
	}

	stratMap := make(map[int]ShardingStrategy)
	for _, s := range strategies {
		stratMap[s.ID] = s
	}

	var calcStrategies []CalculationStrategy
	totalResponseTime := 0.0

	for _, e := range entries {
		strat, ok := stratMap[e.StrategyID]
		if !ok {
			continue
		}
		rt := CalculateResponseTime(strat.AvgResponseTime, e.ShardCount, dataVolume)
		calcStrategies = append(calcStrategies, CalculationStrategy{
			Strategy:   strat,
			ShardCount: e.ShardCount,
			Priority:   e.Priority,
			Comment:    e.Comment,
		})
		totalResponseTime += rt
	}

	avgResponseTime := 0.0
	if len(calcStrategies) > 0 {
		avgResponseTime = totalResponseTime / float64(len(calcStrategies))
	}

	return LoadCalculation{
		ID:                     id,
		Title:                  title,
		Description:            description,
		Status:                 status,
		DataVolume:             dataVolume,
		QueryCount:             queryCount,
		ResultResponseTime:     avgResponseTime,
		ResultLoadDistribution: fmt.Sprintf("%.1f%% на шард", 100.0/float64(len(calcStrategies))),
		Strategies:             calcStrategies,
		StrategyCount:          len(calcStrategies),
	}, nil
}

// GetCalculations возвращает все заявки на расчёт нагрузки.
func (r *Repository) GetCalculations() ([]LoadCalculation, error) {
	entries := []struct {
		StrategyID int
		ShardCount int
		Priority   int
		Comment    string
	}{
		{1, 4, 1, "Для хранения логов по временным диапазонам"},
		{2, 8, 2, "Основная стратегия для пользовательских данных"},
		{3, 3, 3, "Геораспределение для EU и US регионов"},
		{4, 2, 4, "Справочные таблицы с гибким маршрутизацией"},
		{5, 6, 5, "Комбинация Geo + Hash для API-запросов"},
		{6, 10, 6, "Автомасштабирование для пиковых нагрузок"},
	}

	calc, err := r.buildCalculation(
		1,
		"Расчёт нагрузки для e-commerce платформы",
		"Комплексный расчёт распределения нагрузки для интернет-магазина с 50 млн пользователей. Включает анализ всех доступных стратегий шардирования с определением оптимального количества шардов и приоритетов для каждого типа данных.",
		"Рассчитана",
		500,
		10000,
		entries,
	)
	if err != nil {
		return nil, err
	}

	return []LoadCalculation{calc}, nil
}

// GetCalculation возвращает заявку по ID.
func (r *Repository) GetCalculation(id int) (LoadCalculation, error) {
	calcs, err := r.GetCalculations()
	if err != nil {
		return LoadCalculation{}, err
	}

	for _, c := range calcs {
		if c.ID == id {
			return c, nil
		}
	}
	return LoadCalculation{}, fmt.Errorf("заявка не найдена")
}

// GetCalculationForStrategy ищет заявку, содержащую данную стратегию.
func (r *Repository) GetCalculationForStrategy(strategyID int) (*CalculationStrategy, error) {
	calcs, err := r.GetCalculations()
	if err != nil {
		return nil, err
	}

	for _, c := range calcs {
		for _, cs := range c.Strategies {
			if cs.Strategy.ID == strategyID {
				return &cs, nil
			}
		}
	}
	return nil, fmt.Errorf("стратегия не найдена в заявке")
}
