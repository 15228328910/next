package cache

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
)

// Hash 函数用来转换字符串为2的32次方-1中的一个hash位点
type Hash func(data []byte) uint32

// Ring 一致性hash环用来存储节点
type Ring struct {
	// 自定义hash函数
	hash Hash
	// 节点列表
	nodes []int
	// hash环与节点映射
	mp map[int]string
	// 虚拟节点数,防止节点倾斜
	virtual int
}

func NewRing() *Ring {
	return &Ring{
		hash:    crc32.ChecksumIEEE,
		mp:      make(map[int]string, 0),
		virtual: 10,
	}
}

// Add 添加节点
func (r *Ring) Add(node string) {

	for i := 0; i < r.virtual; i++ {
		// 获取节点对应的hash位点
		rand.Seed(int64(i))
		random := rand.Intn(100000)
		nodeName := fmt.Sprintf("%d-%s", random, node)
		index := int(r.hash([]byte(nodeName)))
		r.nodes = append(r.nodes, index)
		sort.Ints(r.nodes)
		r.mp[index] = node
	}
}

// Get 根据key获取数据所在的节点
func (r *Ring) Get(key string) (node string) {
	index := int(r.hash([]byte(key)))
	i := sort.Search(len(r.nodes), func(i int) bool {
		return r.nodes[i] >= index
	})
	position := i % len(r.nodes)
	node = r.mp[r.nodes[position]]
	return
}

// Remove 删除节点
func (r *Ring) Remove(node string) {
	index := int(r.hash([]byte(node)))
	i := sort.SearchInts(r.nodes, index)
	delete(r.mp, r.nodes[i])
	r.nodes = append(r.nodes[:i], r.nodes[i+1:]...)
}
