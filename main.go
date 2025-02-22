package main

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)


func main() {
	fmt.Println("**** PRODUCT-HUNT-GRAPH-VISUALIZE PROJECT ****")

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	chromedp.Run(
		ctx,
		chromedp.ActionFunc(ctx, 
		chromedp.Navigate(""))
	)
}
