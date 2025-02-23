package phuntcrawler

import (
	"context"

	"github.com/chromedp/chromedp"
)

type EntityBuilder struct {
	url       string
	ctx       context.Context
	navigated bool
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
	return ""
}

func (e *EntityBuilder) GetProductDescription() string {
	return ""
}

func (e *EntityBuilder) GetTags() []string {
	return []string{}
}

func (e *EntityBuilder) GetProductTeamMembers() []ProductTeamMember {
	return []ProductTeamMember{}
}

func (e *EntityBuilder) GetPoints() int {
	return 0
}

func (e *EntityBuilder) GetComments() []ProductComments {
	return []ProductComments{}
}

func (e *EntityBuilder) GetDayRank() int {
	return 0
}

func (e *EntityBuilder) GetWeekRank() int {
	return 0
}
