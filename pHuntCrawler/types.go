package phuntcrawler

type IPhuntCrawler interface {
	Crawl() []Product
	GenerateBaseUrls() []string
	ScrapeEntity(entityUrl string) Product
	CollectEntityUrls(baseUrl string) []string
}

type ProductComments struct {
	MemberName    string
	Comment       string
	StarCount     int
	CreatedAtMark string
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
	Comments           []ProductComments
	DayRank            int
	WeekRank           int
}
