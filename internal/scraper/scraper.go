package scraper

import (
	"github.com/anaskhan96/soup"
	"net/http"
	"time"
)

func ExtractJD(url string) (string, error) {
	// Engineering Bias: Always set timeouts. Never wait forever for a website.
	client := &http.Client{Timeout: 10 * time.Second}
	
	resp, err := soup.GetWithClient(url, client)
	if err != nil {
		return "", err
	}
	
	doc := soup.HTMLParse(resp)
	// We look for common job description containers
	mainContent := doc.Find("body")
	if mainContent.Error != nil {
		return "Could not parse body", mainContent.Error
	}
	
	return mainContent.FullText(), nil
}