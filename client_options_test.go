package applesignin

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testTeamID   = "teamID"
	testClientID = "clientID"
	testKeyID    = "keyID"
)

func TestWithBaseURL(t *testing.T) {
	c := New(testTeamID, testClientID, testKeyID, WithBaseURL("https://yet-another-url.com"))
	assertBasicParamsUnchanged(t, c)
	assert.Equal(t, "https://yet-another-url.com", c.baseURL)
}

func TestWithHttpTimeout(t *testing.T) {
	c := New(testTeamID, testClientID, testKeyID, WithHttpTimeout(111*time.Minute))
	assertBasicParamsUnchanged(t, c)
	assert.Equal(t, http.DefaultTransport, c.http.Transport)
	assert.Nil(t, c.http.Jar)
	assert.Nil(t, c.http.CheckRedirect)
	assert.Equal(t, 111*time.Minute, c.http.Timeout)
}

func TestWithHttpTransport(t *testing.T) {
	trans := &http.Transport{}
	c := New(testTeamID, testClientID, testKeyID, WithHttpTransport(trans))
	assertBasicParamsUnchanged(t, c)
	assert.Same(t, trans, c.http.Transport)
	assert.Nil(t, c.http.Jar)
	assert.Nil(t, c.http.CheckRedirect)
	assert.Equal(t, time.Duration(0), c.http.Timeout)
}

func assertBasicParamsUnchanged(t *testing.T, c *Client) {
	assert.Equal(t, testTeamID, c.teamID)
	assert.Equal(t, testClientID, c.clientID)
	assert.Equal(t, testKeyID, c.keyID)
}
