package cassandra

import (
	"fmt"
	"log"
	"reflect"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/gocql/gocql"
)

// Session holds our connection to Cassandra
var Session *gocql.Session

func init() {
	var err error

	appEnv, err := cfenv.Current()
	if err != nil {
		log.Println(err)
		log.Println("Error, cf env not available", appEnv)
		panic("Error, cf env not available")
	}

	cassandraService, err := appEnv.Services.WithName("cassandra-service")

	if err != nil {
		log.Println(err)
		log.Println("Warning, cassandra service not bound", appEnv)
	}

	slice := reflect.ValueOf(cassandraService.Credentials["hostname"])
	c := slice.Len()
	hostNames := make([]interface{}, c)

	for i := 0; i < c; i++ {
		hostNames[i] = slice.Index(i).Interface().(string)
	}

	cluster := gocql.NewCluster(hostNames[0].(string), hostNames[1].(string), hostNames[2].(string))
	user := cassandraService.Credentials["username"].(string)
	password := cassandraService.Credentials["password"].(string)

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: user,
		Password: password,
	}
	cluster.Keyspace = "iotae_dataingestion"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}
