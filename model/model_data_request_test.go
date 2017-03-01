// +build unit
package model_test

import (
	"testing"
	"../model"
	"time"
	"net/url"
	"strconv"
)

const validNodeId = 1

// Test creating a new stor-data request from raw values
func TestNewDataRequest(t *testing.T) {
	nodeId := validNodeId
	start := time.Now()
	end := start.Add(1 * time.Hour)

	dataRequest := model.NewDataRequest(nodeId, start, end)

	checkDataRequestValues(dataRequest, nodeId, start, end, t)
}

// Test creating a new stor-data request from a URL query
func TestFromQuery(t *testing.T) {
	dataRequest := &model.DataRequest{}

	var q = url.Values{}

	nodeId := validNodeId
	start := time.Now()
	end := start.Add(1 * time.Hour)

	q.Add(model.NODE_ID_FIELD, strconv.Itoa(nodeId))
	q.Add(model.START_FIELD, start.Format(time.RFC3339))
	q.Add(model.END_FIELD, end.Format(time.RFC3339))


	dataRequest.FromQuery(q)

	checkDataRequestValues(dataRequest, nodeId, start, end, t)
}

func TestInvalidTimestamp(t *testing.T) {
	dataRequest := &model.DataRequest{}

	var q = url.Values{}

	nodeId := validNodeId
	start := "invalid"
	end := time.Now()

	q.Add(model.NODE_ID_FIELD, strconv.Itoa(nodeId))
	q.Add(model.START_FIELD, start)
	q.Add(model.END_FIELD, end.Format(time.RFC3339))

	if err := dataRequest.FromQuery(q); err == nil {
		t.Error("Expected error")
	}

	q.Set(model.START_FIELD, end.Format(time.RFC3339))
	q.Set(model.END_FIELD, "invalid")

	if err := dataRequest.FromQuery(q); err == nil {
		t.Error("Expected error")
	}

}

// Validates that the values in the struct are the same as expected
func checkDataRequestValues(dataRequest *model.DataRequest, nodeId int, start time.Time, end time.Time, t *testing.T) {
	if dataRequest.NodeId != nodeId {
		t.Errorf("Expected node Id '%d' got '%d'", nodeId, dataRequest.NodeId)
	}

	if dataRequest.Start.Format(time.RFC3339) != start.Format(time.RFC3339) {
		t.Errorf("Expected start '%s' got '%s'",
			start.Format(time.RFC3339),
			dataRequest.Start.Format(time.RFC3339))
	}

	if dataRequest.End.Format(time.RFC3339) != end.Format(time.RFC3339) {
		t.Errorf("Expected start '%s' got '%s'",
			start.Format(time.RFC3339),
			dataRequest.Start.Format(time.RFC3339))
	}
}

func TestValidateStartAfterEndError(t *testing.T) {
	nodeId := validNodeId
	end := time.Now()
	start := end.Add(1 * time.Hour)

	dataRequest := model.NewDataRequest(nodeId, start, end)

	if err := dataRequest.Validate(); err == nil {
		t.Error("Expected validation error, got nothing")
	} else {
		if err.Error() != model.ErrorStartAfterEnd {
			t.Errorf("Expected error \"%s\", got \"%s\"", model.ErrorStartAfterEnd, err)
		}
	}
}