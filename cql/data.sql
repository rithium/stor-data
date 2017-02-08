
CREATE KEYSPACE stor WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  AND durable_writes = true;

CREATE TABLE data (
 nodeId int,
 date timestamp,
 fields map<text, text>,
 PRIMARY KEY (nodeId,date)
) WITH CLUSTERING ORDER BY (date DESC);