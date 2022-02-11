package cache

import (
	"hash/crc32"
	"sort"
)

// Hash 函数用来转换字符串为2的32次方-1中的一个hash位点
type Hash func(data []byte) uint32

// Ring 一致性hash环用来存储节点
type Ring struct {
	hash  Hash
	nodes []int
	mp    map[int]string
}

func NewRing() *Ring {
	return &Ring{
		hash: crc32.ChecksumIEEE,
		mp:   make(map[int]string, 0),
	}
}

// Add 添加节点
func (r *Ring) Add(node string) {
	// 获取节点对应的hash位点
	index := int(r.hash([]byte(node)))
	r.nodes = append(r.nodes, index)
	sort.Ints(r.nodes)
	r.mp[index] = node
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
