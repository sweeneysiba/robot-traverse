package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	Prev    string `json:"prev"`
	Next    string `json:"next"`
	Visited bool   `json:"visited"`
}

func buildNode(a, b, c, d int64) Node {
	return Node{
		Prev:    fmt.Sprintf("%d_%d", a, b),
		Next:    fmt.Sprintf("%d_%d", c, d),
		Visited: false,
	}
}

func getAllPoint(m, n int64) (map[string]Node, map[string]int64) {
	edges := make(map[string]int64)
	var res = make(map[string]Node)
	for i := int64(1); i < n; i++ {
		res[fmt.Sprintf("%d_%d", 0, i)] = buildNode(0, i+1, 0, i-1) // y axis
		res[fmt.Sprintf("%d_%d", m, i)] = buildNode(m, i-1, m, i+1) // x axis
	}
	for i := int64(1); i < m; i++ {
		res[fmt.Sprintf("%d_%d", i, 0)] = buildNode(i-1, 0, i+1, 0) //  y axis
		res[fmt.Sprintf("%d_%d", i, n)] = buildNode(i+1, n, i-1, n) //  x axis

	}
	if n > 0 {
		res[fmt.Sprintf("%d_%d", 0, n)] = buildNode(1, n, 0, n-1)
		edges[fmt.Sprintf("%d_%d", 0, n)] = 1
	}
	if m > 0 {
		res[fmt.Sprintf("%d_%d", m, 0)] = buildNode(m-1, 0, m, 1)
		edges[fmt.Sprintf("%d_%d", 0, n)] = 3
	}
	if n > 0 && m > 0 {
		res[fmt.Sprintf("%d_%d", 0, 0)] = buildNode(0, 1, 1, 0)
		res[fmt.Sprintf("%d_%d", m, n)] = buildNode(m, n-1, m-1, n)
		edges[fmt.Sprintf("%d_%d", 0, 0)] = 0
		edges[fmt.Sprintf("%d_%d", m, n)] = 2
	}

	return res, edges
}

func decideDirection(currentPosition, d string, edges map[string]int64) string {
	switch strings.ToLower(d) {
	case "e":
		return "clock"
	case "w":
		return "anti-clock"
	case "n":
		return "clock"
	case "s":
		return "anti-clock"
	}
	return "invalid direction"
}
func rotate(d, dir string) string {
	switch dir {
	case "l":
		switch d {
		case "e":
			return "n"
		case "w":
			return "s"
		case "n":
			return "w"
		case "s":
			return "e"
		}
	case "r":
		switch d {
		case "e":
			return "s"
		case "w":
			return "n"
		case "n":
			return "e"
		case "s":
			return "w"
		}
	}
	return ""
}
func canMove(pos, d, direction, step string, getAllPoint map[string]Node, edges map[string]int64) (map[string]Node, string, string, string, error) {
	posArr := strings.Split(pos, "_")
	x, _ := strconv.Atoi(posArr[0])
	y, _ := strconv.Atoi(posArr[1])
	// isDirectionModified := false
	nodeVal, exist := getAllPoint[pos]
	if !exist {
		return nil, "", "", "", errors.New("currentPosition is not in the rectangular plane")
	}
	if step == "l" || step == "r" {
		d = rotate(d, step)
		return getAllPoint, pos, d, direction, nil
	}
	value, isedge := edges[fmt.Sprintf("%d_%d", x, y)]
	if isedge {
		switch value {
		case 0:
			if d == "n" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Prev
				direction = "anti-clock"
			} else if d == "e" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Next
				direction = "clock"
			}

		case 1:
			if d == "n" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Next
				direction = "clock"
			} else if d == "w" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Prev
				direction = "anti-clock"

			}
		case 2:
			if d == "w" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Next
				direction = "clock"
			} else if d == "s" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Prev
				direction = "anti-clock"

			}
		case 3:
			if d == "s" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Next
				direction = "clock"

			} else if d == "e" {
				nodeVal.Visited = true
				getAllPoint[pos] = nodeVal
				pos = nodeVal.Prev
				direction = "anti-clock"
			}
		}
	} else if direction == "clock" && !getAllPoint[nodeVal.Next].Visited && step == "m" {
		val := getAllPoint[pos]
		val.Visited = true
		getAllPoint[nodeVal.Next] = val
		pos = nodeVal.Next
	} else if !getAllPoint[nodeVal.Prev].Visited && step == "m" {
		val := getAllPoint[pos]
		val.Visited = true
		getAllPoint[pos] = val
		pos = nodeVal.Prev
	}
	return getAllPoint, pos, d, direction, nil
}

func main() {
	var m, n int64
	var x, y int64
	var d, mov string
	var err error
	var movements []string
	fmt.Println("Give M N.")
	fmt.Scanf("%d %d", &m, &n)
	fmt.Println("Give X Y and Direction(E/W/N/S).")
	fmt.Scanf("%d %d %s", &x, &y, &d)
	fmt.Println("Give Movements(L/R/M).")
	fmt.Scanf("%s", &mov)
	for i := 0; i < len(mov); i++ {
		movements = append(movements, strings.ToLower(string(mov[i])))
	}

	points, edges := getAllPoint(m, n)

	currentPosition := fmt.Sprintf("%d_%d", x, y)
	direction := decideDirection(currentPosition, strings.ToLower(d), edges)
	if direction == "invalid direction" {
		panic(direction)
	}
	DD := decideDirection(currentPosition, d, edges)

	for _, val := range movements {
		points, currentPosition, d, DD, err = canMove(currentPosition, strings.ToLower(d), DD, val, points, edges)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(fmt.Sprintf("current position is (%s, %s) - %s", strings.Split(currentPosition, "_")[0], strings.Split(currentPosition, "_")[1], d))
}

func prettyprint(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
}

/*
STEPS
first input should be (M N) , no space after n and there should be a space between m and n
second input should be (X Y D) , no space after D and there should be a space between X and Y and D
third input should be the moment made by robot no space between any moment ie (MMMRMMLM)

*/
