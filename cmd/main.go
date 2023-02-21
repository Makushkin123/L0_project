package main

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

const (
	portStream = ":4222"
	clusterID  = "test-cluster"
	clientID   = "test"
	portServer = ":4000"
)

func initStream() stan.Conn {

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost"+portStream))

	if err != nil {
		fmt.Println("connect is failed")
	}
	return sc
}

func readFromStream(sc stan.Conn, db *sqlx.DB, cache *Cache) {
	// Simple Async Subscriber
	_, err := sc.Subscribe("foo", func(m *stan.Msg) {

		jsonParsed, err := gabs.ParseJSON(m.Data) //jsonParsed
		fmt.Println(string(m.Data))
		if err != nil {
			fmt.Println("erro format")
		} else {
			key, _ := jsonParsed.Path("order_uid").Data().(string)
			cache.Set(key, m.Data)
			insertDb(key, m, db)
		}
	}, stan.DeliverAllAvailable())

	if err != nil {
		fmt.Println("error receiving messages from the channel")
	}

}

func startServer(cache *Cache) {
	app := &Application{cache: cache}

	srv := &http.Server{
		Addr:    portServer,
		Handler: app.routes(),
	}

	log.Println("Starting the web server on http://127.0.0.1:4000")
	err := srv.ListenAndServe()
	fmt.Println(err)
}

func main() {

	

	db := initDB()
	cache := New(db)
	sc := initStream()


	////Close connection
	defer sc.Close()
	defer db.Close()

	go readFromStream(sc, db, cache)

	startServer(cache)

}
