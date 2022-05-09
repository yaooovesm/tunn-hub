package cache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func testv1() {
	v1 := NewIpTable()
	testCache(v1, "v1")
}

func testv2() {
	v2 := NewIpTableV2(time.Minute*10, time.Minute)
	testCache(v2, "v2")
}

func testv2p() {
	v2p := NewIpTableV2P()
	testCache(v2p, "v2p")
}

func testCacheConcurrent(cache IStrStrCache, name string) {
	start := time.Now().UnixNano()
	wg := sync.WaitGroup{}
	wg.Add(200)
	go func() {
		for i := 0; i < 100; i++ {
			num := i
			go func() {
				for i := 0; i < 5000; i++ {
					cache.Set(strconv.Itoa(i)+strconv.Itoa(num), "1")
				}
				wg.Done()
			}()
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			num := i
			go func() {
				for i := 0; i < 5000; i++ {
					cache.Get(strconv.Itoa(i) + strconv.Itoa(num))
				}
				wg.Done()
			}()
		}
	}()
	wg.Wait()
	end := time.Now().UnixNano()
	fmt.Println(name, " time cost : ", end-start, " ns (", float64(end-start)/1000000, "ms) avg ", float64(end-start)/1000000, " ns")
}

func testCache(cache IStrStrCache, name string) {
	start := time.Now().UnixNano()
	for i := 0; i < 100; i++ {
		num := i
		for i := 0; i < 5000; i++ {
			cache.Set(strconv.Itoa(i)+strconv.Itoa(num), "1")
		}
	}
	for i := 0; i < 100; i++ {
		num := i
		for i := 0; i < 5000; i++ {
			cache.Get(strconv.Itoa(i) + strconv.Itoa(num))
		}
	}
	end := time.Now().UnixNano()
	fmt.Println(name, " time cost : ", end-start, " ns (", float64(end-start)/1000000, "ms) avg ", float64(end-start)/1000000, " ns")
}

func TestIpTable(t *testing.T) {
	for i := 1; i <= 3; i++ {
		fmt.Println()
		fmt.Println("================================== round", i, "start ==================================")
		testv1()
		time.Sleep(time.Second)
		testv2()
		time.Sleep(time.Second)
		testv2p()
		fmt.Println("================================== round", i, "end ==================================")
		fmt.Println()
	}
}
