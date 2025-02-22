package phuntcrawler

import (
	"fmt"
	"slices"
	"time"

	"github.com/fmelihh/product-hunt-graph-visualize/services"
)

const PRODUCT_HUNT_BASE_URL string = "https://www.producthunt.com"

type ServiceDependencies struct {
	baseUrlService   services.BaseUrlService
	entityUrlService services.EntityUrlService
}

type PhuntDomCrawler struct {
	serviceDependencies ServiceDependencies
}

func NewPhuntDomCrawler(serviceDependencies *ServiceDependencies) *PhuntDomCrawler {
	return &PhuntDomCrawler{
		serviceDependencies: *serviceDependencies,
	}
}

func (p *PhuntDomCrawler) crawl() []Product {
	baseUrls := p.generateBaseUrls()
	allEntityUrls := make([]string, 10)

	for _, baseUrl := range baseUrls {
		entityUrls := p.collectEntityUrls(baseUrl)
		for _, entityUrl := range entityUrls {
			p.serviceDependencies.entityUrlService.CreateEntityUrlRecord(entityUrl, baseUrl)
			allEntityUrls = append(allEntityUrls, entityUrl)
		}
		p.serviceDependencies.baseUrlService.CreateBaseUrlRecord(
			baseUrl,
		)
	}
	fmt.Print(allEntityUrls)
	/*
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		for _, url := range baseUrls {
			chromedp.Run(
				ctx,
				chromedp.Navigate(url),
				chromedp.Sleep(2000*time.Millisecond),
			)
		}
		return nil
	*/

	return nil
}

func (p *PhuntDomCrawler) generateBaseUrls() []string {
	baseUrls, err := p.serviceDependencies.baseUrlService.GetAllUrls()
	if err != nil {
		panic(err)
	}

	allPossibleUrls := (func() []string {
		startDate, err := time.Parse("2006/01/02", "2013/12/01")
		if err != nil {
			panic(err)
		}

		today := time.Now()
		var dates []string

		for currentDate := startDate; !currentDate.After(today); currentDate = currentDate.AddDate(0, 0, 1) {
			dates = append(dates, currentDate.Format("2006/01/02"))
		}
		return dates
	})()

	filteredUrls := make([]string, 10, 10)
	for _, url := range allPossibleUrls {
		if !slices.Contains(baseUrls, url) {
			filteredUrls = append(filteredUrls, url)
		}

	}

	return filteredUrls
}

func (p *PhuntDomCrawler) collectEntityUrls(baseUrl string) []string {
	return nil
}

func (p *PhuntDomCrawler) scrapeEntity(entityUrl string) Product {
	return Product{}
}
