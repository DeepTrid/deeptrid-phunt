package services

import (
	"fmt"

	"github.com/fmelihh/product-hunt-graph-visualize/db"
	"github.com/fmelihh/product-hunt-graph-visualize/models"
)

type EntityUrlService struct {
	postgreDb db.PostgreDb
}

func NewEntityUrlService(postgreDb db.PostgreDb) *EntityUrlService {
	return &EntityUrlService{postgreDb: postgreDb}
}

func (h *EntityUrlService) GetAllUrls() ([]string, error) {
	urls := make([]string, 0, 10)

	var entityUrl []models.EntityUrl
	result := h.postgreDb.GormDB.Find(&entityUrl)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find entity url: %v", result.Error)
	}

	fmt.Println("All Entity Urls:")
	for _, entityUrl := range entityUrl {
		urls = append(urls, entityUrl.Url)
	}

	return urls, nil

}

func (h *EntityUrlService) GetEntityUrlByUrl(url string) (models.EntityUrl, error) {
	var entityUrl models.EntityUrl
	err := h.postgreDb.GormDB.Where("url = ?", url).First(&entityUrl).Error
	if err != nil {
		return models.EntityUrl{}, fmt.Errorf("failed to retrieve entity url: %w", err)
	}
	return entityUrl, nil
}

func (h *EntityUrlService) CreateEntityUrlRecord(url string) error {
	entityUrl := models.EntityUrl{Url: url}
	err := h.postgreDb.GormDB.Create(entityUrl).Error
	if err != nil {
		return fmt.Errorf("failed to create entity url: %w", err)
	}
	return nil
}
