package services

import (
	"fmt"

	"github.com/fmelihh/product-hunt-graph-visualize/db"
	"github.com/fmelihh/product-hunt-graph-visualize/models"
)

type BaseUrlService struct {
	postgreDb db.PostgreDb
}

func NewBaseUrlService(postgreDb db.PostgreDb) *BaseUrlService {
	return &BaseUrlService{postgreDb: postgreDb}
}

func (h *BaseUrlService) GetAllUrls() ([]string, error) {
	urls := make([]string, 0, 10)

	var baseUrls []models.BaseUrl
	result := h.postgreDb.GormDB.Find(&baseUrls)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find base url: %v", result.Error)
	}

	fmt.Println("All Base Urls:")
	for _, baseUrl := range baseUrls {
		urls = append(urls, baseUrl.Url)
	}

	return urls, nil
}

func (h *BaseUrlService) GetBaseUrlByUrl(url string) (models.BaseUrl, error) {
	var baseUrl models.BaseUrl
	err := h.postgreDb.GormDB.Where("url = ?", url).First(&baseUrl).Error
	if err != nil {
		return models.BaseUrl{}, fmt.Errorf("failed to retrieve base url: %w", err)
	}
	return baseUrl, nil
}

func (h *BaseUrlService) CreateBaseUrlRecord(url string) error {
	baseUrl := models.NewBaseUrl(url)
	err := h.postgreDb.GormDB.Create(baseUrl).Error
	if err != nil {
		return fmt.Errorf("failed to create base url: %w", err)
	}
	return nil
}
