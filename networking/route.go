package networking

import (
	"fmt"
	log "github.com/cihub/seelog"
	"net"
	"time"
)

//
// IRouter
// @Description:
//
type IRouter interface {
	//
	// Route
	// @Description:
	// @param ip
	// @return interface{}
	//
	Route(ip net.IP) string
	//
	// Add
	// @Description: Add("192.168.1.0/24"[dst cidr],"10.1.1.2"[client])
	// @param dst destination cidr
	// @param tunnel tunnel mark
	// @return error
	//
	Add(dst string, tunnel string) error
	//
	// Delete
	// @Description:
	// @param dst  dst destination cidr
	//
	Delete(dst string)
}

//
// RouteTable
// @Description:
//
type RouteTable struct {
	head            *TableNode
	tail            *TableNode
	Length          int
	rankInterval    int
	rankHalfLife    int
	minimumMaskSize int
}

//
// NewRouteTable
// @Description:
// @return *RouteTable
//
func NewRouteTable(autoRank bool, minimumMaskSize int) *RouteTable {
	tb := &RouteTable{
		Length:          0,
		rankInterval:    3000,
		rankHalfLife:    60000,
		minimumMaskSize: minimumMaskSize,
	}
	if autoRank {
		tb.autoRank()
	}
	return tb
}

//
// Route
// @Description:
// @receiver r
// @param ip
// @return string
//
func (r *RouteTable) Route(ip net.IP) string {
	symbol := ip.To4()
	current := r.head
	for current != nil {
		//匹配依据
		if current.symbol[0] != symbol[0] {
			//move to next
			current = current.Next
			continue
		}
		if current.Net.Contains(ip) {
			current.rank++
			return current.TunnelAddr
		}
		//move to next
		current = current.Next
	}
	return ""
}

//
// autoRank
// @Description:
// @receiver r
//
func (r *RouteTable) autoRank() {
	durationRank := time.Millisecond * time.Duration(r.rankInterval)
	go func() {
		for {
			r.sortByRank()
			time.Sleep(durationRank)
		}
	}()
	durationHalfLife := time.Millisecond * time.Duration(r.rankHalfLife)
	go func() {
		for {
			current := r.head
			for current != nil {
				if current.rank != 0 {
					current.rank /= 2
				}
				current = current.Next
			}
			time.Sleep(durationHalfLife)
		}
	}()
}

//
// sortByRank
// @Description:
// @receiver r
//
func (r *RouteTable) sortByRank() {
	current := r.head
	dummyHead := &TableNode{Next: r.head}
	for current != nil {
		highestRank := findHighestRank(current)
		if highestRank.rank <= current.rank {
			//fmt.Println("skip")
			//fmt.Println("highest --> ", highestRank)
			//fmt.Println("current --> ", current)
			//符合降序则指针前进
			current = current.Next
		} else if highestRank.rank > current.rank {
			//fmt.Println("----------------------------------------------------")
			//fmt.Println("highest --> ", highestRank)
			//fmt.Println("current --> ", current)
			//断开highest
			if highestRank.Last != nil {
				highestRank.Last.Next = highestRank.Next
			}
			if highestRank.Next != nil {
				highestRank.Next.Last = highestRank.Last
			}
			//将highest插入到current之前
			//fmt.Println("current.Last --> ", current.Last)
			if current.Last == nil {
				dummyHead.Next = highestRank
				highestRank.Last = nil
			} else {
				current.Last.Next = highestRank
				current.Last.Next = highestRank
				highestRank.Last = current.Last
			}
			highestRank.Next = current
			current.Last = highestRank
			//fmt.Println()
			//fmt.Println("current.Last --> ", current.Last)
			//fmt.Println("current      --> ", current)
			//fmt.Println("current.Next --> ", current.Next)
			//if current.Next!=nil{
			//	fmt.Println("current.Next.Next --> ", current.Next.Next)
			//}
			//r.head = dummyHead.Next
			//fmt.Println("process --> ", r.ListByName())
			//
			//fmt.Println("----------------------------------------------------")
		}

	}
	r.head = dummyHead.Next
	//fmt.Println()
	//fmt.Println("dummyHead           --> ", dummyHead)
	//fmt.Println("dummyHead.Next      --> ", dummyHead.Next)
	//fmt.Println()
	//fmt.Println("process --> ", r.ListByName())
}

func findHighestRank(head *TableNode) *TableNode {
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return head
	}
	current := head
	max := head
	for current != nil {
		if current.rank > max.rank {
			max = current
		}
		current = current.Next
	}
	return max
}

//
// Add
// @Description: 若已存在则会更新路由表
// @receiver r
// @param dst
// @param uuid
// @return error
//
func (r *RouteTable) Add(dst string, uuid string) error {
	current := r.head
	for current != nil {
		if current.Name == dst {
			if current.TunnelAddr == uuid {
				//存在相同
				return nil
			} else {
				//存在但是目标不同
				log.Info("route updated : ", dst, " --> ", uuid)
				current.TunnelAddr = uuid
				return nil
			}
		}
		//move to next
		current = current.Next
	}
	//不存在，创建并添加
	node, err := NewTableNodeWithLimit(dst, uuid, r.minimumMaskSize)
	if err != nil {
		return err
	}
	if r.head == nil {
		r.head = node
		r.tail = node
	} else {
		r.tail.Next = node
		node.Last = r.tail
		r.tail = node
	}
	r.Length++
	return nil
}

//
// GetByName
// @Description:
// @receiver r
// @param name
// @return *TableNode
//
func (r RouteTable) GetByName(name string) *TableNode {
	current := r.head
	for current != nil {
		if current.Name == name {
			return current
		}
		//move to next
		current = current.Next
	}
	return nil
}

//
// Delete
// @Description:
// @receiver r
// @param dst
//
func (r *RouteTable) Delete(dst string) {
	current := r.head
	for current != nil {
		if current.Net.String() == dst {
			last := current.Last
			next := current.Next
			last.Next = next
			next.Last = last
			r.Length--
			return
		}
		//move to next
		current = current.Next
	}
}

//
// ListByName
// @Description:
// @receiver r
// @return []string
//
func (r RouteTable) ListByName() []string {
	var res []string
	current := r.head
	for current != nil {
		res = append(res, current.Name)
		//fmt.Println("----------------------------------------------")
		//if current.Last!=nil {
		//	fmt.Println("last --> ",current.Last.TunnelAddr)
		//}
		//fmt.Println("current --> ",current.TunnelAddr)
		//if current.Next!=nil {
		//	fmt.Println("next --> ",current.Next.TunnelAddr)
		//}
		//fmt.Println("----------------------------------------------")
		//move to next
		current = current.Next
	}
	return res
}

//
// Print
// @Description:
// @receiver r
//
func (r RouteTable) Print() {
	current := r.head
	for current != nil {
		fmt.Println("----------------------------------------------")
		if current.Last != nil {
			fmt.Println("last --> ", current.Last.TunnelAddr)
		}
		fmt.Println("current --> ", current.TunnelAddr, " (rank=", current.rank, ")")
		if current.Next != nil {
			fmt.Println("next --> ", current.Next.TunnelAddr)
		}
		fmt.Println("----------------------------------------------")
		//move to next
		current = current.Next
	}
}
