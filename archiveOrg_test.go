package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	assert.NoError(t, err)
	assert.Equal(t, uuid, analysis)
}

func TestAnalysisJSON(t *testing.T) {
	analysis, err := testArchiver.analysis(`{"url":"https://example.com","job_id":"spn2-new-job-id"}`)
	assert.NoError(t, err)
	assert.Equal(t, "spn2-new-job-id", analysis)
}

func TestAnalysisMissingJobID(t *testing.T) {
	analysis, err := testArchiver.analysis(`{"url":"https://example.com"}`)
	assert.EqualError(t, err, "archive response contains no job_id")
	assert.Empty(t, analysis)
}

func TestSubmit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "https://github.com/BennyThink", r.Form.Get("url"))
		assert.Equal(t, "on", r.Form.Get("capture_all"))
		_, _ = w.Write([]byte(`{"job_id":"spn2-new-job-id"}`))
	}))
	defer server.Close()

	originalSaveURL := saveUrl
	saveUrl = server.URL
	defer func() { saveUrl = originalSaveURL }()

	response, err := testArchiver.submit("https://github.com/BennyThink")
	assert.NoError(t, err)
	jobID, err := testArchiver.analysis(response)
	assert.NoError(t, err)
	assert.Equal(t, "spn2-new-job-id", jobID)
}
