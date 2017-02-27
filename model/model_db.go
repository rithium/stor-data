package model

import (
	"github.com/gocql/gocql"
	"log"
	"time"
	"../config"
)

type Datastore interface {
	GetLast(int)(*Data, error)
	SaveData(*Data)(error)
	GetData(*DataRequest)([]map[string]interface{}, error)
}

type DB struct {
	*gocql.Session
}

const MAX_RETRIES = 3

func NewDb(params config.CassandraConfig)(*DB, error) {
	cluster := gocql.NewCluster(params.Uri)

	if params.User != "" && params.Pass != "" {
		log.Println("Using Authenticator")

		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: params.User,
			Password: params.Pass,
		}
	}

	cluster.Consistency = gocql.LocalOne
	cluster.ProtoVersion = params.ProtoVersion
	cluster.Keyspace = params.Keyspace

	var session *gocql.Session
	var err error

	session, err = cluster.CreateSession()

	if err != nil {
		var i = 1

		for i <= MAX_RETRIES {
			log.Println("Cassandra retry #", i)

			time.Sleep(2 * time.Second)

			session, err = cluster.CreateSession()

			i++

			if err == nil {
				break
			}
		}

		if err != nil {
			return nil, err
		}

		return &DB{session}, nil
	}

	return &DB{session}, nil
}