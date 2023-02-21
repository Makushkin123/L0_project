package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
)

type Cache struct {
	sync.RWMutex
	items map[string]Item
}

type Item struct {
	Value []byte
}

type Datacahce struct {
	Order_id string
	Data     []byte
}

func New(db *sqlx.DB) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[string]Item)

	//добавляем данные из бд в кэш

	cache := Cache{
		items: items,
	}

	//добавляем данные из бд в кэш
	data := Datacahce{}
	rows, err := db.Queryx("SELECT order_id,data FROM student")

	if err != nil {
		fmt.Println("erro select")
	}
	for rows.Next() {
		err := rows.StructScan(&data)
		if err != nil {
			fmt.Println("erro select and insert in cache")
		}
		cache.Set(data.Order_id, data.Data)
	}

	return &cache
}

func (c *Cache) Set(key string, value []byte) {
	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value: value,
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	// ключ не найден
	if !found {
		return nil, false
	}

	return item.Value, true
}
