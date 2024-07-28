package database

import (
	"fmt"

	"github.com/gocql/gocql"
)

type DatabaseConnection struct {
	Session *gocql.Session
}

var Connection DatabaseConnection

func SetupConnection() {

	cluster := gocql.NewCluster("cassandra:9042")
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.Quorum
	cs, err := cluster.CreateSession()
	Connection.Session = cs
	if err != nil {
		fmt.Println("database error")
		panic(err)
	}
}

func ExecuteQuery(query string, args ...interface{}) error {
	err := Connection.Session.Query(query, args...).Exec()
	return err
}

func SelectQuery(query string, args ...interface{}) *gocql.Query {
	data := Connection.Session.Query(query, args...)
	return data
}

func CheckIfExist(query string, id string) (string, string) {
	var ID string
	var username string
	m := map[string]interface{}{}
	iter := Connection.Session.Query(query, id).Iter()
	for iter.MapScan(m) {
		ID = m["id"].(string)
		username = m["username"].(string)
		m = map[string]interface{}{}
	}
	return ID, username
}
