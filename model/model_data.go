package model

import (
	"time"
	"encoding/json"
	"io"
	"bytes"
	"fmt"
	"log"
	"errors"
	"github.com/gocql/gocql"
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

	return decoder.Decode(d)
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

func (db *DB) GetLast(nodeId int) (*Data, error) {
	var buf bytes.Buffer
	var data Data

	buf.WriteString(fmt.Sprintf("SELECT nodeId, timestamp, fields FROM data WHERE nodeId = ? LIMIT 1"))

	log.Println(db.Query(buf.String(),
		nodeId,
	).String())

	err := db.Query(buf.String(),
		nodeId,
	).Consistency(gocql.One).Scan(&data.NodeId, &data.Timestamp, &data.Data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (db *DB) SaveData(data *Data)(error) {
	if err := data.Validate(); err != nil {
		return err
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("INSERT INTO %s (nodeId, timestamp, fields) values (?, ?, ?) USING TTL ?", "data"))

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