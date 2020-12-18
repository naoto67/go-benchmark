package main

import (
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type M struct {
	Id    int64     `json:"id"`
	Name  string    `json:"string"`
	Flag  bool      `json:"flag"`
	Now   time.Time `json:"now"`
	Array []int     `json:"array"`
}

func main() {
	m := M{Id: 1122334455, Name: "aaa", Flag: true, Now: time.Now(), Array: []int{1, 2, 3, 4, 5}}
	ms := make([]M, 0, 1000)
	for i := 0; i < 1000; i++ {
		ms = append(ms, m)
	}
	var models []M

	b, err := msgpack.Marshal(ms)
	fmt.Println(b, err)
	err = msgpack.Unmarshal(b, &models)
	fmt.Println(err, models)
}
