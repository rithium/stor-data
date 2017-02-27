
CREATE KEYSPACE stor WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  AND durable_writes = true;

CREATE TABLE stor.data (
 nodeId int,
 timestamp timestamp,
 fields map<text, double>,
 PRIMARY KEY (nodeId,timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);