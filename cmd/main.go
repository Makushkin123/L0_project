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
			fmt.Println("неверный формат")
		} else {
			key, _ := jsonParsed.Path("order_uid").Data().(string)
			cache.Set(key, m.Data)
			insertDb(key, m, db)
		}
	}, stan.DeliverAllAvailable())

	if err != nil {
		fmt.Println("ошибка получения сообщений из канала")
	}

}

func startServer(cache *Cache) {
	app := &Application{cache: cache}

	srv := &http.Server{
		Addr:    portServer,
		Handler: app.routes(),
	}

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := srv.ListenAndServe()
	fmt.Println(err)
}

func main() {

	//var objJson string = `{"order_uid": "jinzhu", "age": 18, "tags": ["tag1", "tag2"], "orgs": {"orga": "orga"}}`

	db := initDB()
	defer db.Close()
	cache := New(db)
	sc := initStream()

	//sc.Publish("foo", []byte(objJson))

	////Close connection
	defer sc.Close()
	defer db.Close()

	go readFromStream(sc, db, cache)

	startServer(cache)

}
