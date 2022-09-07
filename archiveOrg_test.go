package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testArchiver = archiveOrg{}

const uuid = "spn2-715fca7db531f22bdf649a2d6450bf63d3effbac"

func TestStatus(t *testing.T) {
	_, err := testArchiver.status(uuid)
	if err != nil {
		assert.Error(t, err, "archive status failed")
	}
}

func TestAnalysis(t *testing.T) {
	file, _ := os.ReadFile("assets/sample.html")
	analysis, err := testArchiver.analysis(string(file))
	if err != nil {
		assert.Error(t, err, "analysis failed")
	}
	assert.Equal(t, uuid, analysis)
}

func TestSubmit(t *testing.T) {
	html, err := testArchiver.submit("https://github.com/BennyThink")
	if err != nil {
		assert.Error(t, err, "archive submit failed")
	}
	assert.Contains(t, html, "spn.watchJob")
}
