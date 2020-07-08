package main

import (
	"bufio"
	"fmt"
	"log"
	// "math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	// "time"
)

type Node struct {
	state [][]string
	x     int
	y     int
	// parent *Node
}

var m int
var n int
var produced = 0
var expanded = 0

var goals [][][]string

var row []int = []int{0, 1, -1, 0}
var col []int = []int{-1, 0, 0, 1}

var checkGoals map[string]bool

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

func convert2DtoString(x [][]string) string {
	arr := make([]string, n*m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			arr[i*n+j] = x[i][j]
		}
	}
	return strings.Join(arr, ",")
}

func dls(src *Node, limit int, depth int) int {
	produced += 1
	if _, ok := checkGoals[convert2DtoString(src.state)]; ok {
		fmt.Println("depth:", depth)
		return 1
	}
	if limit < 0 {
		return -1
	}

	jafarI := src.x
	jafarJ := src.y

	// fmt.Println(jafarI, jafarJ)

	if jafarI+1 >= 0 && jafarI+1 <= (m-1) && jafarJ >= 0 && jafarJ <= (n-1) {
		newNode := new(Node)
		newNode.state = make([][]string, len(src.state))
		for i := range src.state {
			newNode.state[i] = make([]string, len(src.state[i]))
			copy(newNode.state[i], src.state[i])
		}
		temp := newNode.state[jafarI][jafarJ]
		newNode.state[jafarI][jafarJ] = newNode.state[jafarI+1][jafarJ]
		newNode.state[jafarI+1][jafarJ] = temp
		newNode.x = src.x + 1
		newNode.y = src.y
		if dls(newNode, limit-1, depth+1) == 1 {
			fmt.Println("Down")
			return 1
		}
	}

	if jafarI-1 >= 0 && jafarI-1 <= (m-1) && jafarJ >= 0 && jafarJ <= (n-1) {
		newNode := new(Node)
		newNode.state = make([][]string, len(src.state))
		for i := range src.state {
			newNode.state[i] = make([]string, len(src.state[i]))
			copy(newNode.state[i], src.state[i])
		}
		temp := newNode.state[jafarI][jafarJ]
		newNode.state[jafarI][jafarJ] = newNode.state[jafarI-1][jafarJ]
		newNode.state[jafarI-1][jafarJ] = temp
		newNode.x = src.x - 1
		newNode.y = src.y
		if dls(newNode, limit-1, depth+1) == 1 {
			fmt.Println("Up")
			return 1
		}
	}

	if jafarI >= 0 && jafarI <= (m-1) && jafarJ+1 >= 0 && jafarJ+1 <= (n-1) {
		newNode := new(Node)
		newNode.state = make([][]string, len(src.state))
		for i := range src.state {
			newNode.state[i] = make([]string, len(src.state[i]))
			copy(newNode.state[i], src.state[i])
		}
		temp := newNode.state[jafarI][jafarJ]
		newNode.state[jafarI][jafarJ] = newNode.state[jafarI][jafarJ+1]
		newNode.state[jafarI][jafarJ+1] = temp
		newNode.x = src.x
		newNode.y = src.y + 1
		if dls(newNode, limit-1, depth+1) == 1 {
			fmt.Println("Right")
			return 1
		}
	}

	if jafarI >= 0 && jafarI <= (m-1) && jafarJ-1 >= 0 && jafarJ-1 <= (n-1) {
		newNode := new(Node)
		newNode.state = make([][]string, len(src.state))
		for i := range src.state {
			newNode.state[i] = make([]string, len(src.state[i]))
			copy(newNode.state[i], src.state[i])
		}
		temp := newNode.state[jafarI][jafarJ]
		newNode.state[jafarI][jafarJ] = newNode.state[jafarI][jafarJ-1]
		newNode.state[jafarI][jafarJ-1] = temp
		newNode.x = src.x
		newNode.y = src.y - 1
		if dls(newNode, limit-1, depth+1) == 1 {
			fmt.Println("Left")
			return 1
		}
	}

	return -1
}

func ids(src *Node, max_limit int) bool {
	for i := 0; i <= max_limit; i++ {
		if dls(src, i, 0) == 1 {
			fmt.Println("produced:", produced)
			return true
		}
	}
	return false
}

func isSafe(x int, y int) bool {
	return x >= 0 && x < m && y >= 0 && y < n
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

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if start.state[i][j] == "#" {
				start.x, start.y = i, j
			}
		}
	}

	//random actions
	// rand.Seed(time.Now().UnixNano())
	// for i := 0; i < 12; {

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

	checkGoals = make(map[string]bool)
	Permutation(arr2D, 0, m)
	for i := 0; i < len(goals); i++ {
		checkGoals[convert2DtoString(goals[i])] = true
	}
	// for index := range start.state {
	// 	start.state[index] = make([]string, 3)
	// }
	//read from down to up
	fmt.Println(ids(start, 15))

}
