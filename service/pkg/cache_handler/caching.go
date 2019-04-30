package cache_handler

import (
	"fmt"
	"github.com/rajatgpt1521/cachingSystem/service/models"
	"github.com/rs/zerolog/log"
	"sync"
)


type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

// double linked list
type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

// maps string to node in Queue
type Hash map[string]*Node

type cache struct {
	Queue Queue
	Hash  Hash
	lock  sync.RWMutex
}

func NewCache() cache {
	return cache{Queue: NewQueue(), Hash: Hash{}}
}

var (
	//Instance Cache connection object
	Cacheing cache
)

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail}
}

//get data from cache in case not found get from DB
func (c *cache) Read(str string) (err error, string2 string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	node := &Node{}
	if val, ok := c.Hash[str]; ok {
		node = c.remove(val)
		c.addInCache(node)
		c.Hash[str] = node
		log.Info().Msg("got from hash")
		return nil, val.Val
	}

	err, string2 = models.One(str)
	if err != nil {
		log.Error().Err(err).Msg("Data not found in DB")
		return err, "Data not found in DB"
	}

	node = &Node{Val: string2}
	log.Info().Msg("got from DB")
	c.addInCache(node)
	c.Hash[str] = node
	return nil, node.Val

}
func (c *cache) ReadAll() (error, []string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	q := c.Queue
	var str []string
	node := q.Head.Right

	for i := 0; i < q.Length; i++ {

		str = append(str, node.Val)
		node = node.Right
	}

	return nil, str
}

func (c *cache) addInCache(n *Node) {
	//Todo remove print
	fmt.Printf("add: %s\n", n.Val)
	tmp := c.Queue.Head.Right
	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = tmp
	tmp.Left = n

	c.Queue.Length++
	if c.Queue.Length > models.SIZE {
		c.remove(c.Queue.Tail.Left)
	}

}

func (c *cache) remove(n *Node) *Node {

	//Todo remove print
	fmt.Printf("remove: %s\n", n.Val)
	left := n.Left
	right := n.Right
	left.Right = right
	right.Left = left
	c.Queue.Length -= 1

	delete(c.Hash, n.Val)

	return n
}

func (c *cache) Put(str string) (err error, string2 string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	node := &Node{Val: str}
	if err := models.InsertOrUpdate(node.Val); err != nil {
		log.Error().Err(err).Msg("Unable to add in cache DB")
		return err, "Unable to add in cache DB"
	}
	c.addInCache(node)
	c.Hash[str] = node

	return nil, "Successfully added in cache"

}

func (c *cache) display() {
	c.Queue.display()
}

func (q *Queue) display() {
	node := q.Head.Right
	fmt.Printf("%d - [", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.Val)
		if i < q.Length-1 {
			fmt.Printf(" <--> ")
		}
		node = node.Right
	}
	fmt.Println("]")
}

func Initialize() {
	Cacheing = NewCache()
}

func LoadFromDB(){
	if err, cached_data := models.All(); err != nil {
		log.Fatal().Err(err).Msg("Unable to load cache from DB")
	} else {
		for _, data := range cached_data {
			node := &Node{Val: data}
			Cacheing.addInCache(node)

		}

	}
}
func InitializeTestDB() {
	Cacheing = NewCache()
	for _, word := range []string{"cat", "blue", "dog", "tree", "dragon",
		"potato", "house", "tree", "cat"} {
		Cacheing.Put(word)
		Cacheing.display()
		if err, cached_data := models.All(); err != nil {
			log.Fatal().Err(err).Msg("Unable to load cache from DB")
		} else {
			for _, data := range cached_data {
				node := &Node{Val: data}
				Cacheing.addInCache(node)

			}

		}
	}
}

