package phuntcrawler

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

type EntityBuilder struct {
	navigated bool
	url       string
	ctx       context.Context
}

func NewEntityBuilder(url string, ctx context.Context) *EntityBuilder {
	return &EntityBuilder{
		url:       url,
		ctx:       ctx,
		navigated: false,
	}
}

func (e *EntityBuilder) ensureNavigation() error {
	if !e.navigated {
		err := chromedp.Run(e.ctx, chromedp.Navigate(e.url))
		if err != nil {
			return err
		}

		e.navigated = true
	}

	return nil
}

func (e *EntityBuilder) GetProductName() string {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	var title string
	titleXPath := `//div[@class="theme-mirror"]/../div/main/section/div/h1`
	err := chromedp.Run(e.ctx, chromedp.WaitReady(titleXPath), chromedp.Text(titleXPath, &title))
	if err != nil {
		log.Printf("Error getting product name: %v", err)
		return ""
	}
	return title
}

func (e *EntityBuilder) GetProductDescription() string {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	var description string
	descriptionXPath := `//div[@class="theme-mirror"]/../div/main/section[2]`
	err := chromedp.Run(e.ctx, chromedp.WaitReady(descriptionXPath), chromedp.Text(descriptionXPath, &description))
	if err != nil {
		log.Printf("Error getting description: %v", err)
		return ""
	}

	return description
}

func (e *EntityBuilder) GetTags() []string {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	return []string{}
}

func (e *EntityBuilder) GetProductTeamMembers() []ProductTeamMember {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	return []ProductTeamMember{}
}

func (e *EntityBuilder) GetPoints() int {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	return 0
}

func (e *EntityBuilder) GetComments() []ProductComments {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	return []ProductComments{}
}

func (e *EntityBuilder) GetDayRank() int {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	return 0
}

func (e *EntityBuilder) GetWeekRank() int {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	return 0
}
