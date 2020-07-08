package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"github.com/golang-collections/collections/queue"
)

type Node struct {
	parent *Node
	state  [][]string
	x      int
	y      int
}

var expanded int = 0
var produced int = 0
var m int
var n int
var sVisited map[string]*Node
var tVisited map[string]*Node

var row []int = []int{0, 1, -1, 0}
var col []int = []int{-1, 0, 0, 1}

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

func makeNewNode(mat [][]string, x int, y int, newX int, newY int, parent *Node) *Node {
	node := new(Node)
	node.parent = parent
	var cpstate [][]string = make([][]string, len(mat))
	for i := range cpstate {
		cpstate[i] = make([]string, n)
		copy(cpstate[i], mat[i])
	}
	node.state = cpstate
	node.state[x][y], node.state[newX][newY] = node.state[newX][newY], node.state[x][y]
	node.x = newX
	node.y = newY
	return node
}

func isSafe(x int, y int) bool {
	return x >= 0 && x < (m) && y >= 0 && y < (n)
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

func convertStringto2D(x string) [][]string {
	var arr []string = strings.Split(x, ",")
	arr2D := make([][]string, m)
	for index, _ := range arr2D {
		arr2D[index] = make([]string, n)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			arr2D[i][j] = arr[i*n+j]
		}
	}
	return arr2D
}

func isIntersecting() string {
	for key, _ := range sVisited {
		if _, ok := tVisited[key]; ok {
			return key
		}
	}
	return "nil"
}

func printPath(key string) {
	nodeS := sVisited[key]
	nodeT := tVisited[key].parent
	var s []*Node
	for nodeS != nil {
		s = append(s, nodeS)
		nodeS = nodeS.parent
	}

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	// fmt.Println(s)
	for nodeT != nil {
		s = append(s, nodeT)
		nodeT = nodeT.parent
	}

	pValue := s[0]
	depth := 0
	for _, value := range s {
		if value.x > pValue.x {
			fmt.Println("down")
		} else if value.x < pValue.x {
			fmt.Println("up")
		} else if value.y > pValue.y {
			fmt.Println("right")
		} else if value.y < pValue.y {
			fmt.Println("left")
		}
		pValue = value
		depth += 1
	}
	fmt.Print("depth: ", depth-1)
}

func bfs(q *queue.Queue, visited map[string]*Node) {
	temp := queue.New()
	for q.Len() != 0 {
		temp.Enqueue(q.Dequeue().(*Node))
		expanded += 1
	}
	for temp.Len() != 0 {
	current := temp.Dequeue().(*Node)
	for i := 0; i < 4; i++ {
		if isSafe(current.x+row[i], current.y+col[i]) {
			child := makeNewNode(current.state, current.x, current.y, current.x+row[i], current.y+col[i], current)
			_, ok := visited[convert2DtoString(child.state)]
			// fmt.Println(child.state)
			if !ok {
				visited[convert2DtoString(child.state)] = child
				q.Enqueue(child)
				produced += 1
			}
		}
	}
}
	// fmt.Println(visited, "visited")
}
func biDirSearch(src *Node) {
	sQueue := queue.New()
	tQueue := queue.New()
	sVisited[convert2DtoString(src.state)] = src
	sQueue.Enqueue(src)

	for k := 0; k < len(goals); k++ {
		x := 0
		y := 0
		for i := range goals[k] {
			for j := range goals[k][i] {
				if goals[k][i][j] == "#" {
					x = i
					y = j
				}
			}
		}
		node := makeNewNode(goals[k], x, y, x, y, nil)
		tVisited[convert2DtoString(node.state)] = node
		tQueue.Enqueue(node)
	}
	expanded += len(goals)
	for sQueue.Len() != 0 && tQueue.Len() != 0 {
		
		bfs(tQueue, tVisited)
		key0 := isIntersecting()
		if key0 != "nil" {
			fmt.Println("path exists", key0)
			fmt.Println("expanded:", expanded)
			fmt.Println("produced:", produced)
			printPath(key0)
			return
		}
		bfs(sQueue, sVisited)
		// fmt.Println(sVisited)
		// fmt.Println(tVisited)
		key1 := isIntersecting()
		if key1 != "nil" {
			fmt.Println("path exists", key1)
			fmt.Println("expanded:", expanded)
			fmt.Println("produced:", produced)
			printPath(key1)
			return
		}
	
	}

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
	// fmt.Println(start.x, start.y)
	// for index := range start.state {
	// 	start.state[index] = make([]string, 3)
	// }
	sVisited = make(map[string]*Node)
	tVisited = make(map[string]*Node)
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
	// fmt.Println(arr, index)
	var firstPart []string = make([]string, index)
	copy(firstPart, arr[1:index+1])
	var secondPart []string = make([]string, cap(arr)-index)
	copy(secondPart, arr[index+1:cap(arr)])
	arrT := append(firstPart, "#")
	arrSort := append(arrT, secondPart...)
	// fmt.Println(arrSort)

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
	// fmt.Println(goals)
	//read from down to top
	biDirSearch(start)
}
