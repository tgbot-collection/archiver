package main

import "regexp"

const (
	dummyCode = `void(0);`
	zhihu     = `document.getElementsByClassName("Button Modal-closeButton Button--plain")[0].click();`
)

func beautify(url string) string {
	if regexp.MustCompile(`.*zhihu\.com.*`).MatchString(url) {
		return zhihu
	}
	return dummyCode
}
