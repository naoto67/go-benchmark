package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type Model struct {
	Id    int64     `json:"id"`
	Name  string    `json:"string"`
	Flag  bool      `json:"flag"`
	Now   time.Time `json:"now"`
	Array []int     `json:"array"`
}

var data = []byte("hogeharepon!!!")

// シリアライズ系
func Benchmark_JsonMarshalUnmarshalModel(b *testing.B) {
	m := Model{Id: 1122334455, Name: "aaa", Flag: true, Now: time.Now(), Array: []int{1, 2, 3, 4, 5}}
	var model Model
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b, _ := json.Marshal(m)
			json.Unmarshal(b, &model)
		}
	})
	b.StopTimer()
}

func Benchmark_JsonMarshalUnmarshalModels(b *testing.B) {
	m := Model{Id: 1122334455, Name: "aaa", Flag: true, Now: time.Now(), Array: []int{1, 2, 3, 4, 5}}
	ms := make([]Model, 0, 1000)
	for i := 0; i < 1000; i++ {
		ms = append(ms, m)
	}
	var models []Model
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b, _ := json.Marshal(ms)
			json.Unmarshal(b, &models)
		}
	})
	b.StopTimer()
}

func Benchmark_MsgpackModel(b *testing.B) {
	m := Model{Id: 1122334455, Name: "aaa", Flag: true, Now: time.Now(), Array: []int{1, 2, 3, 4, 5}}
	b.ResetTimer()
	var model Model
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b, _ := msgpack.Marshal(m)
			msgpack.Unmarshal(b, &model)
		}
	})
	b.StopTimer()
}

func Benchmark_MsgpackModels(b *testing.B) {
	m := Model{Id: 1122334455, Name: "aaa", Flag: true, Now: time.Now(), Array: []int{1, 2, 3, 4, 5}}
	ms := make([]Model, 0, 1000)
	for i := 0; i < 1000; i++ {
		ms = append(ms, m)
	}
	var models []Model
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b, _ := msgpack.Marshal(ms)
			msgpack.Unmarshal(b, &models)
		}
	})
	b.StopTimer()
}

func Benchmark_MD5(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			md5.Sum(data)
		}
	})
	b.StopTimer()
}

func Benchmark_SHA1(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sha1.Sum(data)
		}
	})
	b.StopTimer()
}

func Benchmark_SHA256(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sha256.Sum256(data)
		}
	})
	b.StopTimer()
}
