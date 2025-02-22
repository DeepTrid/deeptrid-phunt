package main

import (
	"fmt"

	"github.com/fmelihh/product-hunt-graph-visualize/config"
	"github.com/fmelihh/product-hunt-graph-visualize/db"
	phuntcrawler "github.com/fmelihh/product-hunt-graph-visualize/pHuntCrawler"
	"github.com/fmelihh/product-hunt-graph-visualize/services"
)

func main() {
	fmt.Println("**** PRODUCT-HUNT-GRAPH-VISUALIZE PROJECT ****")
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	pdb, err := db.NewPostgreSqlDb(cfg)
	if err != nil {
		panic(err)
	}
	pdb.MigrateDatabaseModels()

	phuntCrawlerServiceDependencies := phuntcrawler.ServiceDependencies{
		BaseUrlService:   *services.NewBaseUrlService(*pdb),
		EntityUrlService: *services.NewEntityUrlService(*pdb),
	}
	pHuntDomCrawler := phuntcrawler.NewPhuntDomCrawler(&phuntCrawlerServiceDependencies)
	pHuntDomCrawler.Crawl()
}
