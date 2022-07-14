package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var hugourl = "https://themes.gohugo.io"

type Theme struct {
	Name string
	Url  string
	Git  string
}

func main() {
	themes := HugoThemeScraper(hugourl)
	file, _ := json.MarshalIndent(themes, "", " ")
	_ = ioutil.WriteFile("hugothemes.json", file, 0644)
}

func HugoThemeScraper(hugourl string) []Theme {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)

	var themedata []Theme

	c.OnHTML(".flex.flex-wrap.justify-left.pr4", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		for _, themelink := range links {
			theme := ScrapeThemePage(themelink)
			if len(theme.Git) > 0 {
				themedata = append(themedata, theme)
			}
		}
	})

	c.Visit(hugourl)
	return themedata
}

func ScrapeThemePage(themeurl string) Theme {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)

	var theme Theme
	theme.Url = themeurl
	splt := strings.Split(themeurl, "/")
	name := splt[len(splt)-2]
	theme.Name = name

	c.OnHTML(".bg-accent-color.br2.hover-bg-primary-color.hover-light-gray.link.ph3.pv2.white", func(e *colly.HTMLElement) {
		gitlink := e.Attr("href")
		if strings.Contains(gitlink, "git") {
			theme.Git = gitlink
		}
	})

	c.Visit(themeurl)
	return theme
}
