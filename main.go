package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	targetfqdn = "github-ranking.com"
)

func IsRelativePath(url string, hostname string) bool {
	if strings.Contains(hostname, url) {
		return true
	} else if strings.Index(url, "//") == 0 {
		return false
	} else if strings.Index(url, "/") == 0 {
		return true
	} else {
		return false
	}
}

func GetAbsoluteURLFromRelativePath(scheme string, fqdn string, relativePath string) string {
	return scheme + "://" + fqdn + relativePath
}

func hasNextPageURL(doc *goquery.Document) (string, bool) {
	nexturl, exists := doc.Find("ul > li.next > a").First().Attr("href")
	if exists == true && IsRelativePath(nexturl, targetfqdn) {
		return GetAbsoluteURLFromRelativePath("https", targetfqdn, nexturl), true
	}
	return nexturl, false
}

func outputRepoAndStar(groupItemSelection *goquery.Selection) {
	groupItemSelection.Each(func(i int, s *goquery.Selection) {
		repositorie := s.Find("span.name span.hidden-lg").Text()
		starts := s.Find("span.stargazers_count").Text()
		fmt.Printf("★ %s : %s\n", strings.TrimSpace(starts), strings.TrimSpace(repositorie))
	})
}

func main() {
	doc, err := goquery.NewDocument(GetAbsoluteURLFromRelativePath("https", targetfqdn, "/repositories"))
	nexturl, hasNext := "", true
	if err != nil {
		log.Fatal(err)
	}
	for hasNext {
		outputRepoAndStar(doc.Find("a.list-group-item"))
		nexturl, hasNext = hasNextPageURL(doc)
		doc, err = goquery.NewDocument(nexturl)
		time.Sleep(3000)
	}
}
