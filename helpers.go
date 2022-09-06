package main

import "regexp"

func beautify(url string) string {

	if regexp.MustCompile(`.*zhihu.com.*`).MatchString(url) {
		return `document.getElementsByClassName("Button Modal-closeButton Button--plain")[0].click();`
	}

	return `void(0);`
}
