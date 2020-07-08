package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	// "math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	// "time"
)

type Node struct {
	parent *Node
	state  [][]string
	x      int
	y      int
	cost   int
	level  int
}

type NodeList []*Node

var explored map[string]*Node
var produced int = 0
var m int
var n int

var row []int = []int{0, 1, -1, 0}
var col []int = []int{-1, 0, 0, 1}

func (lst NodeList) Len() int {
	return len(lst)
}
func (lst NodeList) Less(i, j int) bool {
	return lst[i].cost+lst[i].level < lst[j].cost+lst[j].level
}
func (lst NodeList) Swap(i, j int) {
	lst[i], lst[j] = lst[j], lst[i]
}
func (lst *NodeList) Push(val interface{}) {
	Node := val.(*Node)
	*lst = append(*lst, Node)
}
func (lst *NodeList) Pop() interface{} {
	old := *lst
	n := len(old)
	Node := old[n-1]
	*lst = old[0 : n-1]
	return Node
}

type PQueue struct {
	pq NodeList
}

func NewPQueue() *PQueue {
	queue := new(PQueue)
	queue.pq = make(NodeList, 0)
	heap.Init(&queue.pq)
	return queue
}
func (queue *PQueue) Init() {
	queue.pq = make(NodeList, 0)
	heap.Init(&queue.pq)
}
func (queue *PQueue) Add(value interface{}) {
	heap.Push(&queue.pq, value)
}
func (queue *PQueue) Remove() interface{} {
	return heap.Pop(&queue.pq).(*Node)
}
func (queue *PQueue) Peek() interface{} {
	return queue.pq[0]
}
func (queue *PQueue) Len() int {
	return queue.pq.Len()
}
func (queue *PQueue) IsEmpty() bool {
	return queue.pq.Len() == 0
}

var goals [][][]string

func Permutation(data [][]string, i int, length int) {
	if length == i {
		var cpdata [][]string = make([][]string, length)
		for i := range data {
			cpdata[i] = make([]string, n)
			copy(cpdata[i], data[i])
		}
		goals = append(goals, cpdata)
		return
	}
	for j := i; j < length; j++ {
		swap(data, i, j)
		Permutation(data, i+1, length)
		swap(data, i, j)
	}
}
func swap(data [][]string, x int, y int) {
	data[x], data[y] = data[y], data[x]
}

func makeNewNode(mat [][]string, x int, y int, newX int, newY int, level int, parent *Node) *Node {
	node := new(Node)
	node.parent = parent
	var cpstate [][]string = make([][]string, len(mat))
	for i := range cpstate {
		cpstate[i] = make([]string, n)
		copy(cpstate[i], mat[i])
	}
	node.state = cpstate
	node.state[x][y], node.state[newX][newY] = node.state[newX][newY], node.state[x][y]
	node.cost = math.MaxInt16
	node.level = level
	node.x = newX
	node.y = newY
	return node
}

func calculateCost(initial [][]string) int {
	count := math.MaxInt16
	for k := 0; k < len(goals); k++ {
		temp := 0
		for i := 0; i < len(goals[k]); i++ {
			for j := 0; j < len(goals[k][i]); j++ {
				if initial[i][j] != goals[k][i][j] {
					temp += 1
				}
			}
		}
		if count > temp {
			count = temp
		}
	}
	return count
}

func isSafe(x int, y int) bool {
	return x >= 0 && x < m && y >= 0 && y < n
}

func convert2DtoString(x [][]string) string {
	arr := make([]string, n*m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			arr[i*n+j] = x[i][j]
		}
	}
	return strings.Join(arr, ",")
}

func printPath(node *Node) {
	pnode := node
	depth := 0
	for node != nil {
		depth += 1
		//from target to source
		//read from bottom to up
		if node.x > pnode.x {
			fmt.Println("up")
		} else if node.x < pnode.x {
			fmt.Println("down")
		} else if node.y > pnode.y {
			fmt.Println("left")
		} else if node.y < pnode.y {
			fmt.Println("right")
		}
		pnode = node
		node = node.parent
	}
	fmt.Println(depth)
}
func solve(initial [][]string, x int, y int) int {
	pq := NewPQueue()
	root := makeNewNode(initial, x, y, x, y, 0, nil)
	root.cost = calculateCost(initial)
	expanded := 1
	pq.Add(root)
	explored[convert2DtoString(root.state)] = root
	produced += 1
	for !pq.IsEmpty() {
		min := pq.Remove().(*Node)
		if min.cost == 0 {
			fmt.Println("final state:", min.state)
			fmt.Println("expanded:", expanded)
			fmt.Println("produced:", produced)
			fmt.Println("depth", min.level)
			printPath(min)
			return 1
		}

		for i := 0; i < 4; i++ {
			if isSafe(min.x+row[i], min.y+col[i]) {
				child := makeNewNode(min.state, min.x, min.y, min.x+row[i], min.y+col[i], min.level+1, min)
				_, ok := explored[convert2DtoString(child.state)]
				if !ok {
					child.cost = calculateCost(child.state)
					explored[convert2DtoString(child.state)] = child
					pq.Add(child)
					expanded += 1
				}
			}
		}
		produced += 1
	}
	return -1
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	mn := strings.Split(scanner.Text(), " ")
	mT, _ := strconv.Atoi(mn[0])
	nT, _ := strconv.Atoi(mn[1])
	m, n = mT, nT
	start := new(Node)
	start.state = make([][]string, m)
	for i := 0; i < m; i++ {
		start.state[i] = make([]string, n)
		scanner.Scan()
		start.state[i] = strings.Split(scanner.Text(), " ")
	}

	start.parent = nil

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if start.state[i][j] == "#" {
				start.x, start.y = i, j
			}
		}
	}

	//random actions
	// rand.Seed(time.Now().UnixNano())
	// for i := 0; i < 2; {

	// 	n := rand.Intn(4)

	// 	if isSafe(start.x+row[n], start.y+col[n]) {
	// 		fmt.Println(n)
	// 		i += 1
	// 		start.state[start.x][start.y], start.state[start.x+row[n]][start.y+col[n]] = start.state[start.x+row[n]][start.y+col[n]], start.state[start.x][start.y]
	// 		start.x += row[n]
	// 		start.y += col[n]
	// 	}
	// }
	// var row []int = []int{0, 1, -1, 0}
	// var col []int = []int{-1, 0, 0, 1}
	fmt.Println(start.state)

	explored = make(map[string]*Node)
	arr := make([]string, n*m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			arr[i*n+j] = start.state[i][j]
		}
	}

	sort.SliceStable(arr, func(i, j int) bool {
		if arr[i][len(arr[i])-1] != arr[j][len(arr[j])-1] {
			return arr[i][len(arr[i])-1] < arr[j][len(arr[j])-1]
		} else {
			numI, _ := strconv.Atoi(arr[i][0 : len(arr[i])-1])
			numJ, _ := strconv.Atoi(arr[j][0 : len(arr[j])-1])
			return numI > numJ

		}
	})
	elementP := arr[1]
	count := 1
	index := 0
	flag := false
	for i := 2; i < n*m; i++ {
		element := arr[i]
		if elementP[len(elementP)-1] == element[len(element)-1] {
			count += 1
		} else {
			if count == n {
				count = 1
			} else {
				index = i - n
				flag = true
				break
			}
		}
		elementP = element
	}
	if !flag {
		index = n*m - n
	}

	var firstPart []string = make([]string, index)
	copy(firstPart, arr[1:index+1])
	var secondPart []string = make([]string, cap(arr)-index)
	copy(secondPart, arr[index+1:cap(arr)])
	arrT := append(firstPart, "#")
	arrSort := append(arrT, secondPart...)

	arr2D := make([][]string, m)
	for index, _ := range arr2D {
		arr2D[index] = make([]string, n)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			arr2D[i][j] = arrSort[i*n+j]
		}
	}

	Permutation(arr2D, 0, m)
	//read from up to down
	fmt.Println(solve(start.state, start.x, start.y))
}
