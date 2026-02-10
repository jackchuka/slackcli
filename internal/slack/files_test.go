package slack

import (
	"testing"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestFileFromAPI(t *testing.T) {
	input := slackapi.File{
		ID:                 "F123ABC",
		Name:               "report.pdf",
		Title:              "Quarterly Report",
		Mimetype:           "application/pdf",
		Filetype:           "pdf",
		Size:               102400,
		User:               "U123ABC",
		Created:            slackapi.JSONTime(1700000000),
		URLPrivateDownload: "https://files.slack.com/download/report.pdf",
		Permalink:          "https://team.slack.com/files/report.pdf",
	}

	got := fileFromAPI(input)

	assert.Equal(t, "F123ABC", got.ID)
	assert.Equal(t, "report.pdf", got.Name)
	assert.Equal(t, "Quarterly Report", got.Title)
	assert.Equal(t, "application/pdf", got.Mimetype)
	assert.Equal(t, "pdf", got.Filetype)
	assert.Equal(t, 102400, got.Size)
	assert.Equal(t, "U123ABC", got.User)
	assert.Equal(t, int64(1700000000), got.Created)
	assert.Equal(t, "https://files.slack.com/download/report.pdf", got.URLPrivate)
	assert.Equal(t, "https://team.slack.com/files/report.pdf", got.Permalink)
}

func TestFileFromAPI_Empty(t *testing.T) {
	got := fileFromAPI(slackapi.File{})

	assert.Empty(t, got.ID)
	assert.Empty(t, got.Name)
	assert.Zero(t, got.Size)
	assert.Zero(t, got.Created)
	assert.Empty(t, got.URLPrivate)
}
