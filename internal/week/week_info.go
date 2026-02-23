package week

type WeekInfo struct {
	Week      int
	FirstDate string
	LastDate  string
}

type WeekInfoTemplate struct {
	Week       int
	FirstDate  string
	LastDate   string
	Version    string
	GitHubRepo string
}
