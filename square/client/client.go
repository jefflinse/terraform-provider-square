package client

import (
	"os"

	runtime "github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/google/uuid"
	squareclient "github.com/jefflinse/square-connect/client"
	squaremodel "github.com/jefflinse/square-connect/models"
)

const (
	squareAPIHost = "connect.squareupsandbox.com"
)

// SquareAPI defines an interface for Square's REST API.
type SquareAPI interface {
	DeleteCatalogObject(id string) ([]string, error)
	RetrieveCatalogObject(id string) (*squaremodel.CatalogObject, error)
	UpsertCatalogObject(*squaremodel.CatalogObject) (*squaremodel.CatalogObject, error)
}

// Client is the Square API client.
type Client struct {
	auth   func() runtime.ClientAuthInfoWriter
	square *squareclient.SquareConnect
}

var _ SquareAPI = &Client{}

// NewClient creates a new Square API client using the specified auth token.
func NewClient(token string) *Client {
	transport := httptransport.New(squareAPIHost, squareclient.DefaultBasePath, squareclient.DefaultSchemes)
	if os.Getenv("TERRAFORM_PROVIDER_SQUARE_DEBUG") != "" {
		transport.Debug = true
	}

	squareclient.Default.SetTransport(transport)
	return &Client{
		auth: func() runtime.ClientAuthInfoWriter {
			return httptransport.BearerToken(token)
		},
		square: squareclient.Default,
	}
}

// Generates a new idempotency key for a Square API request.
func newIdempotencyKey() *string {
	key := uuid.New().String()
	return &key
}
