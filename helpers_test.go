package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBeautify(t *testing.T) {

	// assert equality
	var urls = map[string]string{
		"https://www.zhihu.com/question/322578909": zhihu,
		"https://zhuanlan.zhihu.com/p/394575524":   zhihu,
		"https://www.baidu.com/":                   dummyCode,
	}

	for url, code := range urls {
		assert.Equal(t, code, beautify(url), "they should be equal")
	}
}
