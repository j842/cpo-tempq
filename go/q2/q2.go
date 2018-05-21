package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	geo "github.com/paulmach/go.geo"
)

func wktToPts(s string) *geo.Path {
	// Couldn't find a nice WKT parsing library, so do it myself. :/
	// Split string on non-numeric characters, and assign to array.
	// Very little error checking.
	if strings.ToLower(s[0:10]) != "linestring" {
		log.Fatal("Not a LINESTRING: ", s)
	}

	re := regexp.MustCompile("[0-9]+")
	s2 := re.FindAllString(s, -1)

	if len(s2)%2 != 0 {
		log.Fatal("Odd number of coords in LINESTRING!  ", s)
	}

	p := geo.NewPath()

	for i := 0; i+1 < len(s2); i += 2 {
		x, _ := strconv.Atoi(s2[i])
		y, _ := strconv.Atoi(s2[i+1])
		p.Push(geo.NewPoint(float64(x), float64(y)))
	}
	return p
}

type Output struct {
	Does_Intersect bool
	Intersection   []string
}

func roundToInt(x float64) int {
	return int(math.Round(x))
}

func intersect(line1, line2 string) []byte {
	path1 := wktToPts(line1)
	path2 := wktToPts(line2)

	var myout Output

	myout.Does_Intersect = path1.Intersects(path2)

	if myout.Does_Intersect {
		points, _ := path1.Intersection(path2)

		for i, _ := range points {
			// we can only output pointset (WKT MULTIPOINT) with geo, so manually do this as
			// an array of JSON strings, where each string is a point in WKT format.
			s := fmt.Sprintf("POINT(%d %d)", roundToInt(points[i].X()), roundToInt(points[i].Y()))
			myout.Intersection = append(myout.Intersection, s)
		}
	}

	b, err := json.Marshal(myout)

	if err != nil {
		log.Fatal(err)
	}
	return b
}

func tryit(num int, line1, line2 string) {
	fmt.Println("")
	fmt.Println("Example", num)
	fmt.Println(line1, "and", line2)
	r := intersect(line1, line2)
	fmt.Println(string(r))
}

func main() {
	fmt.Println("Question 2")

	tryit(1, "LINESTRING(0 0,2 2)", "LINESTRING(0 2,2 0)")
	tryit(2, "LINESTRING(0 0, 0 2)", "LINESTRING(2 0, 2 2)")
	tryit(3, "LINESTRING(0 0, 0 2, 1 2, 1 0, 2 0, 2 2)", "LINESTRING(0 1, 2 1)")
}

// OUTPUT:
//
// Question 2

// Example 1
// LINESTRING(0 0,2 2) and LINESTRING(0 2,2 0)
// {"Does_Intersect":true,"Intersection":["POINT(1 1)"]}

// Example 2
// LINESTRING(0 0, 0 2) and LINESTRING(2 0, 2 2)
// {"Does_Intersect":false,"Intersection":null}

// Example 3
// LINESTRING(0 0, 0 2, 1 2, 1 0, 2 0, 2 2) and LINESTRING(0 1, 2 1)
// {"Does_Intersect":true,"Intersection":["POINT(0 1)","POINT(1 1)","POINT(2 1)"]}
