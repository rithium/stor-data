package model

import (
	"time"
	"errors"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

type DataRequest struct {
	NodeId	int
	Start	time.Time
	End	time.Time
}

const START_FIELD = "start"
const END_FIELD = "end"

var  (
	ErrorStartAfterEnd = "Start timestamp should be before end"
)

func NewDataRequest(nodeId int, start time.Time, end time.Time) (*DataRequest) {
	dataRequest := DataRequest{}

	dataRequest.NodeId 	= nodeId
	dataRequest.Start 	= start
	dataRequest.End 	= end

	return &dataRequest
}

func (request *DataRequest) FromJson(rawJson *io.Reader) (error) {
	return nil
}

func (request *DataRequest) FromQuery(query url.Values) (error) {
	nodeId, err := strconv.Atoi(query.Get("nodeId"))

	if err != nil {
		return err
	}

	start, err := time.Parse(time.RFC3339, query.Get("start"))

	if err != nil {
		return err
	}

	end, err := time.Parse(time.RFC3339, query.Get("end"))

	if err != nil {
		return err
	}

	request.NodeId 	= nodeId
	request.Start 	= start
	request.End 	= end

	return nil
}

func (request *DataRequest) Validate() (error) {
	if request.Start.IsZero() {
		return errors.New(ErrorInvalidTimestamp)
	}

	if request.End.IsZero() {
		return errors.New(ErrorInvalidTimestamp)
	}

	if request.Start.After(request.End) {
		return errors.New(ErrorStartAfterEnd)
	}

	return nil
}

func (db *DB) GetData(request *DataRequest) ([]map[string]interface{}, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("SELECT * FROM data WHERE nodeId = ? AND date >= ? and date <= ?"))

	result, err := db.Query(buf.String(),
		request.NodeId,
		request.Start,
		request.End,
	).Iter().SliceMap()

	if err != nil {
		return nil, err
	}

	return result, nil
}