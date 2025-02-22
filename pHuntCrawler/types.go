package phuntcrawler

type IPhuntCrawler interface {
	Crawl() []Product
	GenerateBaseUrls() []string
	ScrapeEntity(entityUrl string) Product
	CollectEntityUrls(baseUrl string) []string
}

type ProductSubComments struct {
	MemberName string
	Comment    string
}

type ProductComments struct {
	MemberName  string
	Comment     string
	SubComments []ProductSubComments
}

type ProductTeamMember struct {
	Name     string
	Position string
}

type Product struct {
	ProductName        string
	ProductDescription string
	Tags               []string
	ProductTeamMembers []ProductTeamMember
	Points             int
	Comments           int
	DayRank            int
	WeekRank           int
}
