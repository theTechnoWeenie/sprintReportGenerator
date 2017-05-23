package jira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialsWithoutNewLine(t *testing.T) {
	creds := Credentials{}
	creds.SetCredentials("user\n", "password\\!")
	assert.Equal(t, "user", creds.User)
	assert.Equal(t, "password!", creds.Password)
}

func TestEncodingCredentials(t *testing.T) {
	creds := Credentials{"user", "pass"}
	assert.Equal(t, "dXNlcjpwYXNz", creds.GetEncoded())
}
