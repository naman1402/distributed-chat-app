package database

import (
	"fmt"

	"github.com/gocql/gocql"
)

type DatabaseConnection struct {
	session *gocql.Session
}

var connection DatabaseConnection

func SetupConnection() {

	cluster := gocql.NewCluster("cassandra:9042")
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.Quorum
	cs, err := cluster.CreateSession()
	connection.session = cs
	if err != nil {
		fmt.Println("database error")
		panic(err)
	}
}

func ExecuteQuery() error {
	return nil
}

func SelectQuery() error {
	return nil
}
