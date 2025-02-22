package phuntcrawler

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fmelihh/product-hunt-graph-visualize/services"
)

type PhuntDomCrawler struct {
	historyService services.HistoryService
}

func NewPhuntDomCrawler(historyService *services.HistoryService) *PhuntDomCrawler {
	return &PhuntDomCrawler{
		historyService: *historyService,
	}
}

func (p *PhuntDomCrawler) crawl() []Product {
	urls := p.generateUrls()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	for _, url := range urls {
		chromedp.Run(
			ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(2000*time.Millisecond),
		)
	}
	return nil
}

func (p *PhuntDomCrawler) generateUrls() []string {
	return nil
}
