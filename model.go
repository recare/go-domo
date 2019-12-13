package domo

import (
	"errors"
	"fmt"
	"time"
)

type DataSet struct {
	ID          string    `json:"id,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Schema      *Schema   `json:"schema,omitempty"`
	Rows        int       `json:"rows,omitempty"`
	Columns     int       `json:"columns,omitempty"`
	Owner       *User     `json:"owner,omitempty"`
}

func (d DataSet) validateRequest() error {
	if d.Name == "" {
		return errors.New("missing name")
	}
	if d.Schema == nil {
		return errors.New("missing schema")
	}
	if len(d.Schema.Columns) == 0 {
		return errors.New("len(schema.Columns) is zero")
	}

	return nil
}

type UpdateMethod string

const (
	UpdateMethodAppend  UpdateMethod = "APPEND"
	UpdateMethodReplace UpdateMethod = "REPLACE"
)

type Schema struct {
	Columns Columns `json:"columns,omitempty"`
}

type Column struct {
	Type ColumnType `json:"type,omitempty"`
	Name string     `json:"name,omitempty"`
}

type ColumnType string

const (
	ColumnString   ColumnType = "STRING"
	ColumnDecimal  ColumnType = "DECIMAL"
	ColumnLong     ColumnType = "LONG"
	ColumnDouble   ColumnType = "DOUBLE"
	ColumnDate     ColumnType = "DATE"
	ColumnDateTime ColumnType = "DATETIME"
)

type Columns []Column

type User struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Stream struct {
	ID           int          `json:"id,omitempty"`
	DataSet      *DataSet     `json:"dataSet,omitempty"`
	UpdateMethod UpdateMethod `json:"updateMethod,omitempty"`
}

func (r Stream) validateRequest() error {
	if r.DataSet == nil {
		return errors.New("missing DataSet")
	}
	if r.DataSet.ID != "" {
		return fmt.Errorf("DataSet.ID = %s, expected empty", r.DataSet.ID)
	}
	if err := r.DataSet.validateRequest(); err != nil {
		return fmt.Errorf("invalid DataSet - %w", err)
	}

	if r.UpdateMethod == "" {
		return errors.New("missing UpdateMethod")
	}
	return nil
}

type StreamExecution struct {
	ID           int                  `json:"id,omitempty"`
	StartedAt    time.Time            `json:"startedAt,omitempty"`
	CurrentState StreamExecutionState `json:"currentState,omitempty"`
	CreatedAt    time.Time            `json:"createdAt,omitempty"`
	ModifiedAt   time.Time            `json:"modifiedAt,omitempty"`
}

type StreamExecutionState string

const (
	StreamExecutionActive StreamExecutionState = "ACTIVE"
)

func Time(ts time.Time) string {
	return ts.Format(time.RFC3339)
}
