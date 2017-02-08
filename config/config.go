package config

import (
	"os"
	"strconv"

)

type TomlConfig struct {
	HttpServer	HttpServerConfig
	Cassandra	CassandraConfig
}

type HttpServerConfig struct {
	Uri		string
	Port		int
}

type CassandraConfig struct {
	Uri		string
	User		string
	Pass		string
	Keyspace	string
	ProtoVersion	int
}

var (
	tomlConfig 	TomlConfig
	HttpServer	HttpServerConfig
	Cassandra	CassandraConfig
)

const ENV_HTTP_URL = "S_URL"
const ENV_HTTP_PORT = "S_PORT"

const ENV_CASSANDRA_URL = "S_CASSANDRA_URL"
const ENV_CASSANDRA_USER = "S_CASSANDRA_USER"
const ENV_CASSANDRA_PASS = "S_CASSANDRA_PASS"
const ENV_CASSANDRA_KEYSPACE = "S_CASSANDRA_KEYSPACE"
const ENV_CASSANDRA_PROTO = "S_CASSANDRA_PROTO"

const DEFAULT_HTTP_URL = "0.0.0.0"
const DEFAULT_HTTP_PORT = "80"

const DEFAULT_CASSANDRA_URL = "127.0.0.1"
const DEFAULT_CASSANDRA_USER = ""
const DEFAULT_CASSANDRA_PASS = ""
const DEFAULT_CASSANDRA_KEYSPACE = "stordata"
const DEFAULT_CASSANDRA_PROTO = "4"

// Loads environment variables in to the config
func LoadConfig() {
	tomlConfig.HttpServer.Uri = getEnv(ENV_HTTP_URL, DEFAULT_HTTP_URL)
	tomlConfig.HttpServer.Port, _ = strconv.Atoi(getEnv(ENV_HTTP_PORT, DEFAULT_HTTP_PORT))

	tomlConfig.Cassandra.Uri = getEnv(ENV_CASSANDRA_URL, DEFAULT_CASSANDRA_URL)
	tomlConfig.Cassandra.User = getEnv(ENV_CASSANDRA_USER, DEFAULT_CASSANDRA_USER)
	tomlConfig.Cassandra.Pass = getEnv(ENV_CASSANDRA_PASS, DEFAULT_CASSANDRA_PASS)
	tomlConfig.Cassandra.Keyspace = getEnv(ENV_CASSANDRA_KEYSPACE, DEFAULT_CASSANDRA_KEYSPACE)
	tomlConfig.Cassandra.ProtoVersion, _ = strconv.Atoi(getEnv(ENV_CASSANDRA_PROTO, DEFAULT_CASSANDRA_PROTO))

	HttpServer 	= tomlConfig.HttpServer
	Cassandra	= tomlConfig.Cassandra
}


// Returns the default value if the environment variable is empty
func getEnv(name string, defaultValue string) (string) {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return defaultValue
}