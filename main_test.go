package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	CryptoRand "crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"image"
	_ "image/jpeg"
	"image/png"
	"math/big"
	MathRand "math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/pierrec/lz4"
	"github.com/vmihailenco/msgpack/v5"
)

type Model struct {
	Id    int64     `json:"id" msgpack:"id"`
	Name  string    `json:"string" msgpack:"name"`
	Flag  bool      `json:"flag" msgpack:"flag"`
	Now   time.Time `json:"now" msgpack:"now"`
	Array []int     `json:"array" msgpack:"array"`
}

var (
	data   = []byte("hogeharepon!!!")
	m      = Model{Id: 1122334455, Name: "aaa", Flag: true, Now: time.Now(), Array: []int{1, 2, 3, 4, 5}}
	ms     = make([]Model, 0, 1000)
	msData []byte

	arr = []int{5, 2, 6, 3, 1, 4, 7, 9, 8}
)

func init() {
	for i := 0; i < 1000; i++ {
		ms = append(ms, m)
	}
	msData, _ = json.Marshal(ms)
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func Benchmark_Fib(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fib(20)
		}
	})
	b.StopTimer()
}

func Benchmark_Sort(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var a []int
			copy(a, arr)
			sort.Ints(a)
		}
	})
	b.StopTimer()
}

// シリアライズ系
func Benchmark_JsonMarshalUnmarshalModel(b *testing.B) {
	var model Model
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dat, _ := json.Marshal(m)
			json.Unmarshal(dat, &model)
		}
	})
	b.StopTimer()
}

func Benchmark_JsonMarshalUnmarshalModels(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			models := make([]*Model, 0, 1000)
			dat, _ := json.Marshal(ms)
			json.Unmarshal(dat, &models)
		}
	})
	b.StopTimer()
}

func Benchmark_MsgpackModel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var model Model
			dat, _ := msgpack.Marshal(m)
			msgpack.Unmarshal(dat, &model)
		}
	})
	b.StopTimer()
}

func Benchmark_MsgpackModels(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			models := make([]*Model, 0, 1000)
			dat, _ := msgpack.Marshal(ms)
			msgpack.Unmarshal(dat, &models)
		}
	})
	b.StopTimer()
}

func Benchmark_Gzip(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := new(bytes.Buffer)
			gw, _ := gzip.NewWriterLevel(buf, gzip.BestSpeed)
			gw.Write(msData)
			gw.Close()
			gr, _ := gzip.NewReader(buf)
			gr.Read(buf.Bytes())
		}
	})
	b.StopTimer()
}

func Benchmark_Lz4(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := new(bytes.Buffer)
			lw := lz4.NewWriter(writer)
			lw.Write(msData)
			lw.Close()
			reader := new(bytes.Buffer)
			lr := lz4.NewReader(reader)
			lr.Read(writer.Bytes())
		}
	})
	b.StopTimer()
}

func Benchmark_JpgToPng(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			file, _ := os.Open("./dataset/data1.jpg")
			img, _, _ := image.Decode(file)
			buf := new(bytes.Buffer)
			png.Encode(buf, img)
		}
	})
	b.StopTimer()
}

func Benchmark_AesCfb(b *testing.B) {
	key := []byte("aaaabbbbaaaabbbbaaaabbbbaaaabbbb")
	block, _ := aes.NewCipher(key)
	plainText := []byte("aaaaaaaabbbbaaaa")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cipherText := make([]byte, aes.BlockSize+len(plainText))
			iv := cipherText[:aes.BlockSize]
			iv = []byte("aaaabbbbccccdddd")
			cfb := cipher.NewCFBEncrypter(block, iv)
			ciphertext := make([]byte, len(plainText))
			cfb.XORKeyStream(ciphertext, plainText)

			cfbdec := cipher.NewCFBDecrypter(block, iv)
			plaintextCopy := make([]byte, len(plainText))
			cfbdec.XORKeyStream(plaintextCopy, ciphertext)
		}
	})
	b.StopTimer()
}

func Benchmark_CryptoRandInt10000(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CryptoRand.Int(rand.Reader, big.NewInt(10000))
		}
	})
	b.StopTimer()
}

func Benchmark_MathRandInt10000(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			MathRand.Intn(10000)
		}
	})
	b.StopTimer()
}

func Benchmark_Base64Encoding(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			enc := base64.StdEncoding.EncodeToString(msData)
			base64.StdEncoding.DecodeString(enc)
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
