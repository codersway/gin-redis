package functions

import (
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/util/grand"
)

var (
	wg    sync.WaitGroup
	count = 99999
)

// 耗时11s
// 最原始的添加
func BenchmarkBatchInsert11(b *testing.B) {
	now := time.Now().Unix()

	for i := 0; i <= count; i++ {
		BatchInsert()
	}

	past := time.Now().Unix()
	consuming := past - now
	b.Logf("consuming time: %d", consuming)
}

// 耗时10s
// 10077234371 ns/op
// goroutine添加
func BenchmarkBatchInsert(b *testing.B) {
	now := time.Now().Unix()

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i <= count; i++ {
			BatchInsert()
		}
	}()

	wg.Wait()
	past := time.Now().Unix()

	consuming := past - now
	b.Logf("consuming time: %d", consuming)
}

// 耗时0s到1s
// 0.0996 ns/op
// mset+goroutine
func BenchmarkInsert2(b *testing.B) {
	now := time.Now().Unix()
	wg.Add(1)

	var ss []string
	go func() {
		defer wg.Done()
		for i := 0; i <= count; i++ {

			random := grand.Digits(10)
			ss = append(ss, random)
		}
	}()

	wg.Wait()

	BatchInsert2(ss)

	past := time.Now().Unix()
	consuming := past - now
	b.Logf("consuming time: %d", consuming)
}

func TestDelKeys(t *testing.T) {
	isDel := DelKeys("*78*")

	t.Log(isDel)
}

func TestDelHash(t *testing.T) {
}

func TestBatchInsertSet(t *testing.T) {
	wg.Add(1)

	var ss []string
	go func() {
		defer wg.Done()
		for i := 0; i <= count; i++ {

			random := grand.Digits(10)
			ss = append(ss, random)
		}
	}()

	wg.Wait()

	BatchInsertSet("set", ss)
}

func TestDelSet(t *testing.T) {
	DelSet("set", "*91*")
}

func TestDelZSet(t *testing.T) {
}
