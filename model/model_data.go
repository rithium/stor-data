package model

import (
	"time"
	"encoding/json"
	"io"
	"bytes"
	"fmt"
	"log"
	"errors"
)

const NODE_ID_FIELD = "nodeId"
const TIMESTAMP_FIELD = "timestamp"
const TTL_FIELD = "ttl"
const DATA_FIELD = "data"

var (
	ErrorInvalidTimestamp = "RFC3339 timestamp required"
	ErrorDataFieldMissing = "No data given"
)

type Data struct {
	NodeId		int
	Timestamp	time.Time
	TTL		int
	Data		map[string]float64
}

func NewData(nodeId int, timestamp time.Time, ttl int) (*Data) {
	data := Data{}

	data.NodeId = nodeId
	data.Timestamp = timestamp
	data.TTL = ttl

	return &data
}

func (d *Data) FromJson(rawJson io.Reader) (error) {
	decoder := json.NewDecoder(rawJson)

	if err := decoder.Decode(d); err != nil {
		return err
	}

	return nil
}

func (d *Data) Validate() (error) {
	if d.Timestamp.IsZero() {
		return errors.New(ErrorInvalidTimestamp)
	}

	if len(d.Data) < 1 {
		return errors.New(ErrorDataFieldMissing)
	}

	return nil
}

func (db *DB) SaveData(data *Data)(error) {
	if err := data.Validate(); err != nil {
		return err
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("INSERT INTO %s (nodeId, date, fields) values (?, ?, ?) USING TTL ?", "data"))

	log.Println(db.Query(buf.String(),
		data.NodeId,
		data.Timestamp,
		data.Data,
		data.TTL,
	).String())

	err := db.Query(buf.String(),
		data.NodeId,
		data.Timestamp,
		data.Data,
		data.TTL,
	).Exec()

	return err
}