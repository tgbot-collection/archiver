package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetMD5Hash(t *testing.T) {
	assert.Equal(t, "900150983cd24fb0d6963f7d28e17f72", GetMD5Hash("abc"))
}

func TestTakeScreenshot(t *testing.T) {
	url := "https://github.com/BennyThink"
	filename := takeScreenshot(url)
	expected := fmt.Sprintf("%s.png", GetMD5Hash(url))
	stat, err := os.Stat(expected)
	if err != nil {
		assert.Error(t, err, "file not exist")
	}
	assert.Equal(t, filename, expected)
	assert.Greater(t, stat.Size(), int64(0))
	_ = os.Remove(expected)
}
