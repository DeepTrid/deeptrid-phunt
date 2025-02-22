package phuntcrawler

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fmelihh/product-hunt-graph-visualize/services"
)

const PRODUCT_HUNT_BASE_URL = "https://www.producthunt.com"

type PhuntDomCrawler struct {
	baseUrlService services.BaseUrlService
}

func NewPhuntDomCrawler(baseUrlService *services.BaseUrlService) *PhuntDomCrawler {
	return &PhuntDomCrawler{
		baseUrlService: *baseUrlService,
	}
}

func (p *PhuntDomCrawler) crawl() []Product {
	baseUrls := p.generateBaseUrls()

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
}

func (p *PhuntDomCrawler) generateBaseUrls() []string {
	return nil
}

func (p *PhuntDomCrawler) collectEntityUrls(baseUrl string) []string {
	return nil
}

func (p *PhuntDomCrawler) scrapeEntity(entityUrl string) Product {
	return Product{}
}
