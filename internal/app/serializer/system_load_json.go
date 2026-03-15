package serializer

import (
	"time"

	"web_backend/internal/app/ds"
)

type SystemLoadJSON struct {
	SystemLoadID       uint       `json:"system_load_id"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatorLogin       string     `json:"creator_login"`
	ModeratorLogin     *string    `json:"moderator_login"`
	FormingDate        *time.Time `json:"forming_date"`
	FinishDate         *time.Time `json:"finish_date"`
	Description        *string    `json:"description"`
	CompletedItemCount int        `json:"completed_item_count"`
}

func SystemLoadToJSON(load ds.SystemLoad, creatorLogin, moderatorLogin string, completedItemCount int) SystemLoadJSON {
	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}
	var finishDate *time.Time
	if load.FinishDate.Valid {
		finishDate = &load.FinishDate.Time
	}
	return SystemLoadJSON{
		SystemLoadID:       load.SystemLoadID,
		Status:             load.Status,
		CreatedAt:          load.CreatedAt,
		CreatorLogin:       creatorLogin,
		ModeratorLogin:     mLogin,
		FormingDate:        load.FormingDate,
		FinishDate:         finishDate,
		Description:        load.Description,
		CompletedItemCount: completedItemCount,
	}
}

func SystemLoadFromJSON(j SystemLoadJSON) ds.SystemLoad {
	return ds.SystemLoad{
		Description: j.Description,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
