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

	// creating new cluster
	cluster := gocql.NewCluster("cassandra:9042")
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.Quorum
	// creating a session from the configuration and storing the instance in state variable
	cs, err := cluster.CreateSession()
	Connection.Session = cs
	// error handling
	if err != nil {
		fmt.Println("database error")
		panic(err)
	}
}

// state variables stores the session instance, this function takes the query and args
// and pass it through the session and exec it, returns error
func ExecuteQuery(query string, args ...interface{}) error {
	err := Connection.Session.Query(query, args...).Exec()
	return err
}

// creates Query from query and args, and returns it
func SelectQuery(query string, args ...interface{}) *gocql.Query {
	data := Connection.Session.Query(query, args...)
	return data
}

// creating mapping from string to interface{}
// passing query and id as Query into session, and turn result into iterator(allows you to go through the results one row at a time, without loading everything into memory at once)
// scan through every row of iterator, populates it with a row that is returned from cassandra.
// In each iteration, it extracts the id and username from the current row, assigns them to variables (ID and username), and then resets the map m for the next row.
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
