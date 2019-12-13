package domo

import (
	"net/http"
	"time"
)

const (
	baseUrl = "https://api.domo.com"
	version = "v1"

	pathAuth    = "oauth/token"
	pathDataset = "datasets"
	pathStream  = "streams"
)

type Client interface {
	GetDataSet(id string) (*DataSet, error)
	CreateDataSet(*DataSet) (*DataSet, error)
	UpdateDataSet(*DataSet) (*DataSet, error)
	Import(dataSetID string, updateMethod UpdateMethod, fields [][]string) error
	CreateStream(*Stream) (*Stream, error)
	CreateStreamExecution(streamID int) (*StreamExecution, error)
	StreamImport(streamID, executionID, part int, data [][]string) error
	CommitStreamExecution(streamID, executionID int) error
}

type HttpClient struct {
	clientID, secret string
	client           *http.Client
	token            *token
}

func New(clientID, secret string) *HttpClient {
	return &HttpClient{
		clientID: clientID,
		secret:   secret,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HttpClient) SetHttpClient(client *http.Client) {
	c.client = client
}
