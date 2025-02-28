package phuntcrawler

import (
	"context"
	"log"
	"slices"
	"strings"

	"github.com/chromedp/cdproto/cdp"
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
	var tags []string
	var nodes []*cdp.Node

	tagsXpath := `//div[@data-sentry-component="TagList"]/a`
	err := chromedp.Run(e.ctx, chromedp.WaitReady(tagsXpath), chromedp.Nodes(tagsXpath, &nodes))
	if err != nil {
		log.Printf("Error getting description: %v", err)
		return []string{}
	}

	for _, node := range nodes {
		if node.NodeType == 1 {
			nodeValue := strings.TrimSpace(node.Children[0].NodeValue)
			if !slices.Contains(tags, nodeValue) {
				tags = append(tags, strings.TrimSpace(node.Children[0].NodeValue))
			} else {
				break
			}
		}
	}

	return tags
}

func (e *EntityBuilder) GetProductTeamMembers() []ProductTeamMember {
	if err := e.ensureNavigation(); err != nil {
		panic(err)
	}
	var nodes []*cdp.Node
	var teamMembers []ProductTeamMember

	teamMembersButton := `//a[text() = "Team"]`
	teamMembersXpath := `//div[@data-sentry-component="InfiniteScroll"]/section[contains(@data-test,"maker-card")]`

	err := chromedp.Run(e.ctx,
		chromedp.WaitReady(teamMembersButton),
		chromedp.Click(teamMembersButton),
		chromedp.ActionFunc(
			func(ctx context.Context) error {
				return chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil).Do(ctx)
			}),
		chromedp.WaitReady(teamMembersXpath),
		chromedp.Nodes(teamMembersXpath, &nodes),
	)

	if err != nil {
		log.Printf("Error getting description: %v", err)
		return []ProductTeamMember{}
	}

	for _, node := range nodes {
		var name, position string
		err := chromedp.Run(e.ctx,
			chromedp.Text(node.FullXPath()+`/div[2]/a[1]`, &name),
			chromedp.Text(node.FullXPath()+`/div[2]/a[2]`, &position),
		)
		if err != nil {
			log.Println("Error extracting data:", err)
			continue
		}
		teamMembers = append(teamMembers, ProductTeamMember{
			Name:     name,
			Position: position,
		})
	}

	return teamMembers
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
