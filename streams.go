package domo

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c HttpClient) CreateStream(in *Stream) (*Stream, error) {
	if err := in.validateRequest(); err != nil {
		return nil, fmt.Errorf("invalid CreateStream input - %w", err)
	}
	data, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("error marshalling create stream body - %w", err)
	}

	body, err := c.do(http.MethodPost, pathStream, &payload{data: bytes.NewReader(data), contentType: "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error executing create stream request - %w", err)
	}

	var out *Stream
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body of create stream request - %s", err.Error())
	}

	return out, nil
}

func (c HttpClient) CreateStreamExecution(streamID int) (*StreamExecution, error) {
	body, err := c.do(http.MethodPost,
		fmt.Sprintf("%s/%d/executions", pathStream, streamID), nil)
	if err != nil {
		return nil, fmt.Errorf("error executing create stream execution request - %w", err)
	}

	var out *StreamExecution
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body of create stream execution request - %s", err.Error())
	}

	return out, nil
}

func (c HttpClient) StreamImport(streamID, executionID, part int, fields [][]string) error {
	b := bytes.NewBuffer(nil)
	w := csv.NewWriter(b)

	if err := w.WriteAll(fields); err != nil {
		return fmt.Errorf("error csv marshalling csv - %w", err)
	}

	if _, err := c.do(http.MethodPut,
		fmt.Sprintf("%s/%d/executions/%d/part/%d", pathStream, streamID, executionID, part),
		&payload{data: b, contentType: "text/csv"}); err != nil {
		return fmt.Errorf("error executing stream import request - %w", err)
	}

	return nil
}

func (c HttpClient) CommitStreamExecution(streamID, executionID int) error {
	if _, err := c.do(http.MethodPut,
		fmt.Sprintf("%s/%d/executions/%d/commit", pathStream, streamID, executionID), nil); err != nil {
		return fmt.Errorf("error executing commit stream request - %w", err)
	}

	return nil
}
