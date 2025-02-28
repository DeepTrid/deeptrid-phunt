package phuntcrawler

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func createChromadpContext(ctx context.Context) (context.Context, context.CancelFunc) {
	chromedpContext, cancelChromeDp := chromedp.NewContext(ctx)
	return chromedpContext, cancelChromeDp
}

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

func (e *EntityBuilder) ensureNavigation(externalContext context.Context) error {
	if !e.navigated {
		var err error
		if externalContext == context.TODO() {
			err = chromedp.Run(e.ctx, chromedp.Navigate(e.url))
		} else {
			err = chromedp.Run(externalContext, chromedp.Navigate(e.url))
		}

		if err != nil {
			return err
		}

		e.navigated = true
	}

	return nil
}

func (e *EntityBuilder) GetProductName() string {
	if err := e.ensureNavigation(context.TODO()); err != nil {
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
	if err := e.ensureNavigation(context.TODO()); err != nil {
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
	if err := e.ensureNavigation(context.TODO()); err != nil {
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
	if err := e.ensureNavigation(context.TODO()); err != nil {
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
	if err := e.ensureNavigation(context.TODO()); err != nil {
		panic(err)
	}

	var nodes []*cdp.Node
	starXpath := `//div[@data-sentry-component="StarRating"]/label`
	err := chromedp.Run(e.ctx,
		chromedp.WaitReady(starXpath),
		chromedp.Nodes(starXpath, &nodes),
	)
	if err != nil {
		log.Printf("Error getting description: %v", err)
		return 0
	}

	starCount := 0
	for _, node := range nodes {
		if node.NodeType == 1 {
			starCount += 1
		}
	}
	return starCount
}

func (e *EntityBuilder) GetComments() []ProductComments {
	if err := e.ensureNavigation(context.TODO()); err != nil {
		panic(err)
	}
	reviewsButtonXpath := `//a[text() = "Reviews"]`
	showMoreButtonXpath := `//section[@data-sentry-component="ReviewsFeedWithLoadButton"]//button[starts-with(text(),"Show")]`

	err := chromedp.Run(e.ctx,
		chromedp.WaitReady(reviewsButtonXpath),
		chromedp.Click(reviewsButtonXpath),
	)

	if err != nil {
		log.Println("Error checking button visibility:", err)
		return []ProductComments{}
	}

	var currentHeight int64 = 0
	var previousHeight int64 = 0

	for {
		var nodes []*cdp.Node

		time.Sleep(2 * time.Second)
		err := chromedp.Run(e.ctx,
			chromedp.Nodes(showMoreButtonXpath, &nodes, chromedp.AtLeast(0)),
		)

		if err != nil {
			log.Println("Error checking button visibility:", err)
			break
		}

		if len(nodes) > 0 {
			err = chromedp.Run(e.ctx,
				chromedp.Click(showMoreButtonXpath),
			)
			if err != nil {
				log.Println("Error clicking button:", err)
				break
			}
		} else {
			break
		}

		err = chromedp.Run(e.ctx,
			chromedp.Evaluate(`document.body.scrollHeight`, &currentHeight),
		)
		if err != nil {
			log.Println("Error getting scroll height:", err)
			break
		}

		err = chromedp.Run(e.ctx, chromedp.ActionFunc(
			func(ctx context.Context) error {
				return chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil).Do(ctx)
			}))

		if err != nil {
			log.Println("Error scrolling page:", err)
			break
		}

		if currentHeight <= previousHeight {
			fmt.Println("Reached end of scroll.")
			break
		}

		previousHeight = currentHeight
		time.Sleep(1 * time.Second)

		fmt.Println("Scrolled comments section...")
	}

	var commentNodes []*cdp.Node
	var comments []ProductComments

	commentsXpath := `//div[@data-sentry-component="ReviewItemWithComments"]`
	err = chromedp.Run(e.ctx,
		chromedp.WaitReady(commentsXpath),
		chromedp.Nodes(commentsXpath, &commentNodes),
	)

	if err != nil {
		log.Printf("Error getting comment: %v", err)
		return []ProductComments{}
	}

	for _, node := range commentNodes {
		ratingCount := 0
		ratingXpath := node.FullXPath() + `//div[@data-sentry-component="ReviewItemWithComments"]//div[@data-sentry-component="StarRating"]`
		var ratingNodes []*cdp.Node

		err := chromedp.Run(
			e.ctx,
			chromedp.WaitReady(ratingXpath),
			chromedp.Nodes(ratingXpath, &ratingNodes),
		)
		if err == nil {
			for _, node := range ratingNodes {
				if node.NodeType == 1 {
					ratingCount += 1
				}
			}
		}

		var memberName, commentString, createdAtMark string

		err = chromedp.Run(e.ctx,
			chromedp.Text(node.FullXPath()+`//div[@data-sentry-component="UserSpotlight"]/div/a`, &memberName),
			chromedp.Text(node.FullXPath()+`//div[@data-sentry-component="ReviewItemWithComments"]/div/div[2]/div[3]`, &commentString),
			chromedp.AttributeValue(node.FullXPath()+`//div[@data-sentry-component="ReviewItemWithComments"]/div/div[2]/div[4]/div/time`, "datetime", &createdAtMark, nil),
		)
		if err != nil {
			log.Println("Error extracting data:", err)
			continue
		}
		comments = append(comments, ProductComments{
			MemberName:    memberName,
			Comment:       commentString,
			StarCount:     ratingCount,
			CreatedAtMark: createdAtMark,
		})
	}

	return comments
}

func (e *EntityBuilder) GetDayRank() int {
	if err := e.ensureNavigation(context.TODO()); err != nil {
		panic(err)
	}
	return 0
}

func (e *EntityBuilder) GetWeekRank() int {
	if err := e.ensureNavigation(context.TODO()); err != nil {
		panic(err)
	}
	return 0
}
