## go-domo
Go API client for domo. 

This client is incomplete, but implements select methods pertaining to data imports from the DataSet and Stream APIs.

Currently:
```go
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
```

## Install
`go get -u github.com/veyo-care/go-domo`

## Usage

For instance, to create a dataset and perform an import using the DataSet api:
```go
import "github.com/veyo-care/go-domo"

client := domo.New("client-id", "secret")

ds, err := client.CreateDataSet(&domo.DataSet{
	Name:        "DataSet Name",
	Description: "Description",
	Schema: &domo.Schema{
		Columns: domo.Columns{
			{Name: "String", Type: domo.ColumnString},
			{Name: "Number", Type: domo.ColumnDecimal},
		},
	},
})
if err != nil {
	panic(err)
}

if err := client.Import(ds.ID,
	domo.UpdateMethodAppend,
	[][]string{{"A", "1"}, {"B", "2"}}); err != nil {
	panic(err)
}
```

## API references

https://developer.domo.com/docs/dataset-api-reference/dataset

https://developer.domo.com/docs/streams-api-reference/streams
