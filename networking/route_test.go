package networking

import (
	"fmt"
	log "github.com/cihub/seelog"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"testing"
	"time"
)

func TestRouting(t *testing.T) {
	table := RouteTable{}
	err := table.Add("192.168.1.0/24", "1")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = table.Add("192.168.2.0/24", "2")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = table.Add("192.168.3.0/24", "3")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(table.ListByName())
	t1 := table.Route(net.ParseIP("192.168.2.2"))
	fmt.Println("192.168.2.2 --> ", t1)
	//fmt.Println(table.ListByName())

	err = table.Add("192.168.4.0/24", "4")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(table.ListByName())
	t1 = table.Route(net.ParseIP("192.168.4.2"))
	fmt.Println("192.168.4.2 --> ", t1)

	for i := 0; i < 10; i++ {
		t1 = table.Route(net.ParseIP("192.168.3.6"))
	}
	table.Print()
	fmt.Println()
	fmt.Println()
	table.sortByRank()
	table.Print()
	fmt.Println()
	fmt.Println()
	for i := 0; i < 50; i++ {
		t1 = table.Route(net.ParseIP("192.168.4.6"))
	}
	fmt.Println()
	fmt.Println()
	for i := 0; i < 100; i++ {
		table.sortByRank()
	}
	table.Print()
	fmt.Println()
	fmt.Println()
	//fmt.Println(table.ListByName())

}
func TestRouting2(t *testing.T) {
	go func() {
		log.Info("pprof serve on 0.0.0.0:9988")
		if err := http.ListenAndServe("0.0.0.0:9988", nil); err != nil {
			_ = log.Warn("pprof start failed : ", err)
		}
	}()
	s1 := 20
	s2 := 50
	table := NewRouteTable(true, 8)
	for i := 0; i < s1; i++ {
		for j := 0; j < s2; j++ {
			cidr := "10." + strconv.Itoa(i) + "." + strconv.Itoa(j) + ".0/16"
			//fmt.Println(cidr)
			err := table.Add(cidr, strconv.Itoa(i)+"."+strconv.Itoa(j))
			if err != nil {
				fmt.Println("err --> ", err)
				return
			}
		}
	}
	fmt.Println("table size -> ", table.Length)
	time.Sleep(time.Minute * 10)
}

func TestRouting1(t *testing.T) {
	go func() {
		log.Info("pprof serve on 0.0.0.0:9988")
		if err := http.ListenAndServe("0.0.0.0:9988", nil); err != nil {
			_ = log.Warn("pprof start failed : ", err)
		}
	}()
	size := 100000000
	s1 := 10
	s2 := 10
	table := NewRouteTable(true, 8)
	for i := 0; i < s1; i++ {
		for j := 0; j < s2; j++ {
			cidr := "10." + strconv.Itoa(i) + "." + strconv.Itoa(j) + ".0/16"
			//fmt.Println(cidr)
			err := table.Add(cidr, strconv.Itoa(i)+"."+strconv.Itoa(j))
			if err != nil {
				fmt.Println("err --> ", err)
				return
			}
		}
	}
	fmt.Println("table size -> ", table.Length)
	fmt.Println("sample size -> ", size)
	var testList []net.IP
	for i := 0; i < size; i++ {
		x := rand.Intn(s1)
		y := rand.Intn(s2)
		z := rand.Intn(s1-1) + 1
		addr := "10." + strconv.Itoa(x) + "." + strconv.Itoa(y) + "." + strconv.Itoa(z)
		testList = append(testList, net.ParseIP(addr))
	}
	//var testList1 []net.IP
	//for i := 0; i < size; i++ {
	//	x := rand.Intn(254)
	//	y := rand.Intn(254)
	//	z := rand.Intn(253) + 1
	//	addr := "10." + strconv.Itoa(x) + "." + strconv.Itoa(y) + "." + strconv.Itoa(z)
	//	testList1 = append(testList1, net.ParseIP(addr))
	//}
	//go func() {
	//	for {
	//		time.Sleep(time.Millisecond * 100)
	//		table.sortByRank()
	//	}
	//}()
	//fmt.Println("random test size --> ", len(testList))
	start := time.Now().UnixNano()
	//for i := range testList {
	//	//res := table.Route(testList[i])
	//	//fmt.Println("[", i+1, "]res -> ", res)
	//	table.Route(testList[i])
	//}
	for i := 0; i < size; i++ {
		x := rand.Intn(s1)
		y := rand.Intn(s2)
		z := rand.Intn(s1-1) + 1
		addr := "10." + strconv.Itoa(x) + "." + strconv.Itoa(y) + "." + strconv.Itoa(z)
		table.Route(net.ParseIP(addr))
	}
	fmt.Println()
	fmt.Println()
	end := time.Now().UnixNano()
	fmt.Println("---------------------------------------")
	fmt.Println("random networking test")
	total := end - start
	fmt.Println("total = ", total, "ns (", float64(total)/1000000, "ms)")
	fmt.Println("avg = ", float64(total)/float64(size), "ns (", float64(total)/float64(size)/1000000, "ms)")

	str := strconv.FormatFloat(1000/(float64(total)/float64(size)/1000000), 'f', -1, 64)
	fmt.Println("speed(avg) = ", str, "r/s")
	fmt.Println("---------------------------------------")

	//table.sortByRank()
	fmt.Println()
	//fmt.Println("rank...")
	fmt.Println()
	start = time.Now().UnixNano()
	//for i := 0; i < size; i++ {
	//	x := rand.Intn(254)
	//	y := rand.Intn(254)
	//	z := rand.Intn(253) + 1
	//	addr := "10." + strconv.Itoa(x) + "." + strconv.Itoa(y) + "." + strconv.Itoa(z)
	//	table.Route(net.ParseIP(addr))
	//}
	for i := range testList {
		//res := table.Route(testList[i])
		//fmt.Println("[", i+1, "]res -> ", res)
		table.Route(testList[i])
	}
	end = time.Now().UnixNano()
	fmt.Println("---------------------------------------")
	fmt.Println("random networking test")
	total = end - start
	fmt.Println("total = ", total, "ns (", float64(total)/1000000, "ms)")
	fmt.Println("avg = ", float64(total)/float64(size), "ns (", float64(total)/float64(size)/1000000, "ms)")
	str = strconv.FormatFloat(1000/(float64(total)/float64(size)/1000000), 'f', -1, 64)
	fmt.Println("speed(avg) = ", str, "r/s")
	fmt.Println("---------------------------------------")
}

func TestRouting3(t *testing.T) {
	table := NewRouteTable(true, 8)
	err := table.Add("172.23.0.2/32", "1")
	if err != nil {
		fmt.Println("err -> ", err)
		return
	}
	res := table.Route(net.ParseIP("172.23.0.2"))
	fmt.Println(res)
}

func TestNet(t *testing.T) {
	ip, ipNet, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ip)
	fmt.Println(ipNet.Mask.Size())
}
