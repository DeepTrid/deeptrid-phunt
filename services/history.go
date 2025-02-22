package services

import (
	"fmt"

	"github.com/fmelihh/product-hunt-graph-visualize/db"
	"github.com/fmelihh/product-hunt-graph-visualize/models"
)

type HistoryService struct {
	postgreDb db.PostgreDb
}

func NewHistoryService(postgreDb db.PostgreDb) *HistoryService {
	return &HistoryService{postgreDb: postgreDb}
}

func (h *HistoryService) GetAllUrls() ([]string, error) {
	urls := make([]string, 0, 10)

	var histories []models.History
	result := h.postgreDb.GormDB.Find(&histories)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find histories: %v", result.Error)
	}

	fmt.Println("All Histories:")
	for _, history := range histories {
		urls = append(urls, history.Url)
	}

	return urls, nil

}

func (h *HistoryService) GetHistoryByUrl(url string) (models.History, error) {
	var history models.History
	err := h.postgreDb.GormDB.Where("url = ?", url).First(&history).Error
	if err != nil {
		return models.History{}, fmt.Errorf("failed to retrieve history: %w", err)
	}
	return history, nil
}

func (h *HistoryService) CreateHistoryRecord(url string) error {
	history := models.History{Url: url}
	err := h.postgreDb.GormDB.Create(history).Error
	if err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}
	return nil
}
