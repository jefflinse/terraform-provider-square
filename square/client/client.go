package client

import (
	"os"

	runtime "github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	squareclient "github.com/jefflinse/square-connect/client"
)

const (
	squareAPIHost = "connect.squareupsandbox.com"
)

// Square is the Square API client.
type Square struct {
	auth   func() runtime.ClientAuthInfoWriter
	square *squareclient.SquareConnect
}

// NewClient creates a new Square API client using the specified auth token.
func NewClient(token string) *Square {
	transport := httptransport.New(squareAPIHost, squareclient.DefaultBasePath, squareclient.DefaultSchemes)
	if os.Getenv("TERRAFORM_PROVIDER_SQUARE_DEBUG") != "" {
		transport.Debug = true
	}

	squareclient.Default.SetTransport(transport)
	return &Square{
		auth: func() runtime.ClientAuthInfoWriter {
			return httptransport.BearerToken(token)
		},
		square: squareclient.Default,
	}
}
