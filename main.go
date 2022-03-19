package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Graph structure
type Graph struct {
	rooms []*Room
}

//Room structure
type Room struct {
	key      string
	adjacent []*Room
	path     []*Room
	visited  bool
	occupied bool
}

// Reads file and returns a string slice
func readAntsFile(filename string) []string {
	file, _ := os.Open(filename)
	scanned := bufio.NewScanner(file)
	scanned.Split(bufio.ScanLines)

	var lines []string

	for scanned.Scan() {
		lines = append(lines, scanned.Text())
	}

	return lines
}

func NumAnts(s []string) int {
	antNum := s[0]
	s = readAntsFile("ants.txt")
	if s[0] <= "0" {
		err := fmt.Errorf("invalid number of ants")
		fmt.Println(err.Error())
	}

	antNumInt, _ := strconv.Atoi(antNum)
	return antNumInt
}

// Gets out the start room and returns it
func StartRoom([]string) string {

	var startRoom string
	s := readAntsFile("ants.txt")

	for i := 0; i < len(s); i++ {
		if s[i] == "##start" {
			startRoom = strings.Split(string(s[i+1]), " ")[0]
		}

	}

	return startRoom

}

// Gets out the end room and returns it
func EndRoom([]string) string {
	var endRoom string
	s := readAntsFile("ants.txt")

	for i := 0; i < len(s); i++ {
		if s[i] == "##end" {
			endRoom = strings.Split(string(s[i+1]), " ")[0]
		}

	}
	//fmt.Println(endRoom)
	return endRoom
}

var (
	StartR = StartRoom(readAntsFile("ants.txt"))
	EndR   = EndRoom(readAntsFile("ants.txt"))
)

//Add Room to a graph
func (g *Graph) AddRoom(k string) {
	if contains(g.rooms, k) {
		err := fmt.Errorf("Room %v not added because it is an existing key", k)
		fmt.Println(err.Error())
	} else {
		g.rooms = append(g.rooms, &Room{key: k})
	}
}

//getRoom returns a pointer to the Room key integer
func (g *Graph) getRoom(k string) *Room {
	for i, v := range g.rooms {
		if v.key == k {
			return g.rooms[i]
		}
	}
	return nil
}

//contains checks if the Room key exists
func contains(s []*Room, k string) bool {
	for _, v := range s {
		if k == v.key {
			return true
		}
	}
	return false
}

func doesContain(s string, sl []string) bool {
	for _, word := range sl {
		if s == word {
			return true
		}
	}
	return false
}

func doesContainRoom(sl []*Room, s string) bool {

	for _, word := range sl {
		if s == word.key {
			return true
		}
	}
	return false
}

func (g *Graph) PrintPath() {
	fmt.Println(StartRoom(readAntsFile("ants.txt")))
	for _, v := range g.rooms {
		for _, r := range v.path {
			fmt.Println(r)
		}
	}
}

// add all edges
// Add edge to the graph. deals with a directional graph only but condition in the main makes it undirected
func (g *Graph) AddEdge(from, to string) {
	//get Room
	fromRoom := g.getRoom(from)
	toRoom := g.getRoom(to)

	//check error
	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("invalid edge (%v-->%v)", from, to)
		fmt.Println(err.Error())
	} else if contains(fromRoom.adjacent, to) {
		err := fmt.Errorf("existing edge (%v-->%v)", from, to)
		fmt.Println(err.Error())
	} else if fromRoom == toRoom {
		err := fmt.Errorf("cannot connect room to itself (%v --> %v)", from, to)
		fmt.Println(err.Error())
	} else if fromRoom.key == EndR {
		//toRoom.adjacent = append(toRoom.adjacent, fromRoom)
		//} //else if toRoom.key == StartR {
		//toRoom.adjacent = append(toRoom.adjacent, fromRoom)
	} else {
		fromRoom.adjacent = append(fromRoom.adjacent, toRoom)
	}
	//add edge etc
}

//Print will print the adjacent list for each Room of the graph
func (g *Graph) Print() {
	// fmt.Println(readAntsFile("ants.txt"))
	fmt.Printf("The number of ants is: %v ", NumAnts(readAntsFile("ants.txt")))
	fmt.Println()

	for _, v := range g.rooms {
		if v.key == StartR {
			fmt.Printf("\n Start Room is %v : ", StartR)

		} else if v.key == EndR {
			fmt.Printf("\n End Room is %v :", EndR)

		} else {
			fmt.Printf("\n Room %v : ", v.key)

		}
		for _, v := range v.adjacent {
			fmt.Printf(" %v,", v.key)
		}
	}
	fmt.Println()
}

// global variable which will store all of the valid paths in a slice of slices of string.
var validPaths [][]*Room

// Slices to hold paths from both algos for comparison
var dfsPaths [][]*Room
var bfsPaths [][]*Room

// Depth First Search algorithm that operates recursively
func DFS(r *Room, g Graph) {

	vList := []string{}
	sRoom := g.getRoom(StartR)

	// set the room being checked visited status to true
	if r.key != EndR {
		r.visited = true

		// append the r key to the visited list
		vList = append(vList, r.key)

		// range through the neighbours of the r
		for _, nbr := range r.adjacent {
			if !nbr.visited {
				/* for each neighbour that hasn't been visited,
				- append their key to the visited slice,
				- then apply dfs to them recursively,
				- then append their key to their path value
				*/

				//fmt.Println("*", vList)
				nbr.path = append(r.path, nbr)
				if contains(nbr.path, EndR) {
					// for _, r := range nbr.path {
					// 	fmt.Println("DFS Finale: ----->", r.key)
					// }
					dfsPaths = append(dfsPaths, nbr.path)

				}
				//fmt.Println(nbr.path)
				vList = append(vList, nbr.key)
				//fmt.Println("La Final1: ", vList)

				//DeleteEdge(r, g)
				DFS(nbr, Graph{g.rooms})

			}

		}

	} else {

		if len(sRoom.adjacent) > 1 && !contains(sRoom.adjacent, EndR) {
			vList = append(vList, r.key)

			sRoom.adjacent = sRoom.adjacent[1:][:]

			DFS(sRoom, Graph{g.rooms})

		} else {
			vList = append(vList, r.key)

		}
	}

}

// Depth first search function that operates recursively
func DFSBFS(r *Room, g Graph) bool {

	//vList := []string{}
	//sRoom := g.getRoom(StartR)

	// set the room being checked visited status to true
	if r.key != EndR {
		r.visited = true

		// range through the neighbours of the r
		for _, nbr := range r.adjacent {
			if !nbr.visited {
				/* for each neighbour that hasn't been visited,
				- append their key to the visited slice,
				- then apply dfs to them recursively,
				- then append their key to their path value
				*/

				//fmt.Println("*", vList)
				nbr.path = append(r.path, nbr)
				if contains(nbr.path, EndR) {
					//fmt.Println("dfs print check", nbr.path)

					return true

				}

			}

		}

	}
	return false
}

// Function to remove room using its key
func RemoveRoomIndex(s []*Room, index string) []*Room {
	ret := make([]*Room, 0)

	for i := 0; i < len(s); i++ {
		if s[i].key == index {
			ret = append(ret, s[:i]...)
			ret = append(ret, s[i+1:]...)
		}
	}
	return ret
}

// function to remove an element from a slice of ints
func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

//removes a string from a slice (unordered)
func remove(s []*Room, k string) []*Room {
	for i := 0; i < len(s); i++ {
		if s[i].key == k {
			s[i] = s[len(s)-1]

		}

	}
	return s[:len(s)-1]
}

//delete edge from starting room
func DeleteEdge(r *Room, g Graph) {

	for i := 0; i < len(r.path); i++ {
		for _, room := range g.rooms {
			//	for _ , edge := range room.adjacent
			for j := 0; j < len(room.adjacent); j++ {
				if room.adjacent[j] == r.path[i] {
					room.adjacent = remove(room.adjacent, r.key)
				}
			}
		}
	}
}

// Breadth First Search algorithm
func BFS(r *Room, g Graph) {

	var vPaths [][]*Room

	//queue variable, procedurally populated with rooms yet to be visited
	var queue []*Room

	//set start room as visited
	r.visited = true

	//initialise queue with start room
	queue = append(queue, r)

	// }

	// checks if there is a link between start and end directly
	for i, v := range g.getRoom(StartR).adjacent {
		if v.key == g.getRoom(EndR).key {
			g.getRoom(EndR).path = append(g.getRoom(EndR).path, g.getRoom(StartR))
			vPaths = append(vPaths, g.getRoom(EndR).path)
			g.getRoom(StartR).adjacent = append(g.getRoom(StartR).adjacent[:i], g.getRoom(StartR).adjacent[i+1:]...)
			//fmt.Println("End reached --------------------------------------------------------------------->", g.getRoom(StartR).path)
		}

	}

	//checks the queue for the end room and if the queue is not empty

	for !contains(queue, g.getRoom(EndR).key) && len(queue) >= 1 {
		qfront := queue[0]
		//	fmt.Println("QF:", qfront.key)

		for _, room := range qfront.adjacent {
			if !room.visited {
				room.visited = true
				room.path = append(qfront.path, room)
				//
				queue = append(queue, room)
				// fmt.Println(queue)
			}

		}

		queue = queue[1:]

		if doesContainRoom(queue, g.getRoom(EndR).key) {

			for _, room := range g.rooms {
				room.visited = false
			}
			vPaths = append(vPaths, qfront.path)

			DeleteEdge(qfront, g)

			if len(g.getRoom(StartR).adjacent) == 0 {
				break
			}
			//g.Print()

			if len(g.getRoom(StartR).adjacent) >= 1 {
				for _, froom := range g.getRoom(StartR).adjacent {
					for _, sroom := range froom.adjacent {
						if sroom.key != g.getRoom(EndR).key {
							//	fmt.Println("loop 8")
							break
						} else {
							//	fmt.Println("loop 9")
							BFS(g.getRoom(StartR), Graph{g.rooms})
							queue = queue[1:]
						}
					}
				}
			}
			BFS(g.getRoom(StartR), Graph{g.rooms})

		}
	}
	for _, v := range vPaths {
		v = append(v, g.getRoom(EndR))
		bfsPaths = append(bfsPaths, v)
		//fmt.Printf("BFS Finale: ----> ")
		// for _, r := range v {
		// 	fmt.Printf("%v ", r.key)
		// }
		fmt.Println()
	}

}

func PathSelection(bfs [][]*Room, dfs [][]*Room) [][]*Room {
	// a = bfsPaths
	// b = dfsPaths

	bfsPathNum := len(bfs)
	dfsPathNum := len(dfs)

	if bfsPathNum > dfsPathNum {
		validPaths = append(validPaths, bfsPaths...)
	} else if dfsPathNum > bfsPathNum {
		validPaths = append(validPaths, dfsPaths...)
	} else if dfsPathNum == bfsPathNum {

		bfscounter := 0
		dfscounter := 0

		for _, path := range bfs {

			bfscounter += len(path)

		}

		for _, path := range dfs {
			dfscounter += len(path)
		}

		if bfscounter > dfscounter {
			validPaths = append(validPaths, bfs...)
		} else if dfscounter > bfscounter {
			validPaths = append(validPaths, dfs...)
		} else if bfscounter == dfscounter {
			validPaths = append(validPaths, bfs...)
		}

	}
	return validPaths

}

//error checking for duplicated starting points
// func PathDupeCheck(path [][]*Room) [][]*Room {

// 	var result [][]*Room
// 	//	var compare [][]*Room
// 	var output [][]*Room

// 	if len(path) == 1 {
// 		return path
// 	}

// 	if len(path) == 2 {
// 		if path[0][0].key != path[1][0].key {
// 			return path
// 		} else if len(path[0]) <= len(path[1]) {
// 			output = append(output, path[0])
// 			return output
// 		}
// 		output = append(output, path[1])
// 		return output
// 	}

// 	if len(path) > 2 {

// 		dataMap:=  make(map[*Room][]*Room)

// 		for _ , slice := range path{
// 			dataMap[slice[0]] = slice
// 		}

// 	}
// }

func main() {

	bfsGraph := Graph{}

	//adding all rooms
	for i, line := range readAntsFile("ants.txt") {
		if strings.Contains(string(line), " ") {
			bfsGraph.AddRoom(strings.Split(readAntsFile("ants.txt")[i], " ")[0])
		}
		// adding all edges from and to rooms

		if strings.Contains(string(line), "-") {
			bfsGraph.AddEdge(strings.Split(readAntsFile("ants.txt")[i], "-")[0], strings.Split(readAntsFile("ants.txt")[i], "-")[1])
			bfsGraph.AddEdge(strings.Split(readAntsFile("ants.txt")[i], "-")[1], strings.Split(readAntsFile("ants.txt")[i], "-")[0])
		}

	}

	//bfsGraph.Print()
	BFS(bfsGraph.getRoom(StartR), bfsGraph)
	//bfsGraph.Print()

	dfsGraph := Graph{}

	//adding all rooms
	for i, line := range readAntsFile("ants.txt") {
		if strings.Contains(string(line), " ") {
			dfsGraph.AddRoom(strings.Split(readAntsFile("ants.txt")[i], " ")[0])
		}
		// adding all edges from and to rooms
		// only adding edges in one direction to make the graph directional
		if strings.Contains(string(line), "-") {
			dfsGraph.AddEdge(strings.Split(readAntsFile("ants.txt")[i], "-")[0], strings.Split(readAntsFile("ants.txt")[i], "-")[1])
		}

	}

	//dfsGraph.Print()
	DFS(dfsGraph.getRoom(StartR), dfsGraph)

	// for _, path := range PathDupeCheck(PathSelection(bfsPaths,dfsPaths)) {
	// 	for _, room := range path {
	// 		fmt.Print(room.key)
	// 		fmt.Print("  ")
	// 	}
	// 	fmt.Println()
	// }

	// for _, path := range PathSelection(bfsPaths, dfsPaths)[0] {

	// 	fmt.Print(path.key)

	// }
	// fmt.Println()

	// 	for _, path := range PathDupeCheck(PathSelection(bfsPaths, dfsPaths)) {
	// 		for _, room := range path {
	// 	fmt.Print(room.key)
	// }
	// fmt.Println()
	// 	}

	// fmt.Println(PathSelection(bfsPaths,dfsPaths)[0])
	// fmt.Println(PathSelection(bfsPaths,dfsPaths))

	// for _, paths := range PathSelection(bfsPaths,dfsPaths){
	// fmt.Println(paths)
	// }

	//fmt.Println(PathDupeCheck(PathSelection(bfsPaths, dfsPaths)))
	//fmt.Println(len(PathSelection(bfsPaths, dfsPaths)[0][0].key))
	//fmt.Printf("%T", PathSelection(bfsPaths, dfsPaths)[0][0])

	//fmt.Println(PathSelection(bfsPaths, dfsPaths)[1])

	dataMap := make(map[*Room][]*Room)

	//var compare [][]*Room  = PathSelection(bfsPaths, dfsPaths)[1:]

	//allKeys := make(map[string]bool)

	for _, item := range PathSelection(bfsPaths, dfsPaths) {
		if value, ok := dataMap[item[0]]; !ok {
			dataMap[item[0]] = item
		} else {
			if len(item) <= len(value) {
				dataMap[item[0]] = item

			}
		}
	}
for key , value := range dataMap{
	fmt.Print(key.key)  
	for _ , room := range value{
	fmt.Print(room.key)
}
fmt.Println()
}


}
