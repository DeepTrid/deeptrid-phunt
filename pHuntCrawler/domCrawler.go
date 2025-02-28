package phuntcrawler

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/fmelihh/product-hunt-graph-visualize/services"
)

const PRODUCT_HUNT_BASE_URL string = "https://www.producthunt.com"

type ServiceDependencies struct {
	BaseUrlService   services.BaseUrlService
	EntityUrlService services.EntityUrlService
}

type PhuntDomCrawler struct {
	serviceDependencies ServiceDependencies
}

func NewPhuntDomCrawler(serviceDependencies *ServiceDependencies) *PhuntDomCrawler {
	return &PhuntDomCrawler{
		serviceDependencies: *serviceDependencies,
	}
}

func (p *PhuntDomCrawler) Crawl() []Product {
	baseUrls := p.GenerateBaseUrls()
	allEntityUrls := make([]string, 10)

	for _, baseUrl := range baseUrls {
		entityUrls := p.CollectEntityUrls(baseUrl)
		for _, entityUrl := range entityUrls {
			p.serviceDependencies.EntityUrlService.CreateEntityUrlRecord(entityUrl, baseUrl)
			allEntityUrls = append(allEntityUrls, entityUrl)
			fmt.Printf("%s entity url was successfulyy added to db", entityUrl)
		}
		p.serviceDependencies.BaseUrlService.CreateBaseUrlRecord(
			baseUrl,
		)
		fmt.Printf("%s base url was successfully added to db, total entity url is %d", baseUrl, len(allEntityUrls))
		time.Sleep(30 * time.Second)
	}

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

func (p *PhuntDomCrawler) GenerateBaseUrls() []string {
	baseUrls, err := p.serviceDependencies.BaseUrlService.GetAllUrls()
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

	filteredUrls := make([]string, 0, 10)
	for _, url := range allPossibleUrls {
		url = PRODUCT_HUNT_BASE_URL + "/leaderboard/daily/" + url
		if !slices.Contains(baseUrls, url) {
			filteredUrls = append(filteredUrls, url)
		}
	}
	return filteredUrls
}

func (p *PhuntDomCrawler) CollectEntityUrls(baseUrl string) []string {
	resp, err := http.Get(baseUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("an error occurred while %s sending. status code %d", baseUrl, resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	var productURLs []string

	doc.Find("section").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Find("a").Attr("href")
		if exists {
			productURLs = append(productURLs, fmt.Sprintf("%s%s", PRODUCT_HUNT_BASE_URL, href))
		}
	})

	return productURLs
}

func (p *PhuntDomCrawler) ScrapeEntity(entityUrl string) Product {
	contextValue, cancel := createChromadpContext(context.Background())
	defer cancel()

	entityBuilder := NewEntityBuilder(entityUrl, contextValue)
	comments := entityBuilder.GetComments()
	product := Product{
		ProductName:        entityBuilder.GetProductName(),
		ProductDescription: entityBuilder.GetProductDescription(),
		Tags:               entityBuilder.GetTags(),
		ProductTeamMembers: entityBuilder.GetProductTeamMembers(),
		Points:             entityBuilder.GetPoints(),
		Comments:           comments,
		DayRank:            entityBuilder.GetDayRank(),
		WeekRank:           entityBuilder.GetWeekRank(),
	}
	return product
}
