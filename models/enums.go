package models

type Status string

const (
	StatusCrawled    Status = "CRAWLED"
	StatusNotCrawled Status = "NOT_CRAWLED"
	StatusError      Status = "ERROR"
)
