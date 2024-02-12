// Package di have all the injections dependency logic
package di

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionProvider interface for Session methods.
type SessionProvider interface {
	Session() (client.ConfigProvider, error)
}

// SessionConfig struct for Session configuration.
type SessionConfig struct {
	CredentialsFile string
	Endpoint        string
	Region          string `env:"AWS_REGION" envDefault:"us-east-1"`
}

// Session attributes required for SessionProvider.
type Session struct {
	session *session.Session
	config  *SessionConfig
}

// Session method for create session client
func (s *Session) Session() (client.ConfigProvider, error) {
	return s.session, nil
}

// newSessionProvider instantiate new SessionProvider.
func newSessionProvider(config *SessionConfig) SessionProvider {
	return &Session{
		config: config,
	}
}

// newAWSSessionProvider provider to aws session
func newAWSSessionProvider() SessionProvider {
	return newSessionProvider(&SessionConfig{})
}
