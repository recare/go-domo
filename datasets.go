package domo

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c HttpClient) CreateDataSet(in *DataSet) (*DataSet, error) {
	if err := in.validateRequest(); err != nil {
		return nil, fmt.Errorf("invalid CreateDataSet input - %w", err)
	}
	data, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("error marshalling CreateDataSet body - %w", err)
	}

	body, err := c.do(http.MethodPost, pathDataset, &payload{data: bytes.NewReader(data), contentType: "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error performing CreateDataSet request - %w", err)
	}

	var out *DataSet
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body of CreateDataSet request - %w", err)
	}

	return out, nil
}

func (c HttpClient) UpdateDataSet(in *DataSet) (*DataSet, error) {
	if err := in.validateRequest(); err != nil {
		return nil, fmt.Errorf("invalid UpdateDataSet input - %w", err)
	}
	data, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("error marshalling UpdateDataSet body - %w", err)
	}

	body, err := c.do(http.MethodPut, fmt.Sprintf("%s/%s", pathDataset, in.ID), &payload{data: bytes.NewReader(data), contentType: "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error performing UpdateDataSet request - %w", err)
	}

	var out *DataSet
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body of UpdateDataSet request - %w", err)
	}

	return out, nil
}

func (c HttpClient) GetDataSet(id string) (*DataSet, error) {
	body, err := c.do(http.MethodGet, url(pathDataset, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error performing GetDataSet request - %w", err)
	}

	var out *DataSet
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body of GetDataSet request - %w", err)
	}

	return out, nil
}

func (c HttpClient) Import(dataSetID string, updateMethod UpdateMethod, fields [][]string) error {
	b := bytes.NewBuffer(nil)
	w := csv.NewWriter(b)

	if err := w.WriteAll(fields); err != nil {
		return fmt.Errorf("error csv marshalling csv - %w", err)
	}

	_, err := c.do(http.MethodPut,
		fmt.Sprintf("%s/%s/data?updateMethod=%s", pathDataset, dataSetID, updateMethod),
		&payload{data: b, contentType: "text/csv"})
	if err != nil {
		return fmt.Errorf("error performing Import request - %w", err)
	}

	return nil
}
