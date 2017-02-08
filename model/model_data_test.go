// +build unit
package model_test

import (
	"testing"
	"../model"
	"encoding/json"
	"bytes"
	"time"
)

func TestNewDataFromJson(t *testing.T) {
	data := model.Data{}

	nodeId := 1
	timestamp := time.Now()
	ttl := 0
	var dataMap map[string]interface{}

	rawJson := map[string]interface{}{
		model.NODE_ID_FIELD: nodeId,
		model.TIMESTAMP_FIELD: timestamp.Format(time.RFC3339),
		model.TTL_FIELD: ttl,
		model.DATA_FIELD: dataMap,
	}

	b, _ := json.Marshal(rawJson)

	reader := bytes.NewBuffer(b)

	if err := data.FromJson(reader); err != nil {
		t.Error("Expected time error")
	}

	if data.NodeId != nodeId {
		t.Error("Expected node id %d got %d", nodeId, data.NodeId)
	}

	if data.Timestamp.Format(time.RFC3339) != timestamp.Format(time.RFC3339) {
		t.Error("Expected timestamp to be in RFC3339 format", nodeId, data.NodeId)
	}

	if data.TTL != ttl {
		t.Error("Expected ttl %d got %d")
	}
}

func TestNewDataFromJsonInvalidTimestamp(t *testing.T) {
	data := model.Data{}

	rawJson := map[string]interface{}{
		model.NODE_ID_FIELD: "1",
		model.TIMESTAMP_FIELD: "1",
		model.TTL_FIELD: "1",
		model.DATA_FIELD: map[string]interface{}{
			"1": "1",
		},
	}

	b, _ := json.Marshal(rawJson)

	reader := bytes.NewBuffer(b)

	if err := data.FromJson(reader); err == nil {
		t.Error("Expected time error")
	}
}

func TestPostData(t *testing.T) {

}