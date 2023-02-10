package sendfigures1

import "fmt"

type Cube struct {
	Length int
	Width  int
	Height int
	Area   int
	Volume int
}

func Process() {
	cubes := []Cube{{1, 2, 5, 0, 0}, {4, 6, 7, 0, 0}, {8, 10, 2, 0, 0}, {7, 2, 7, 0, 0}, {6, 1, 8, 0, 0}}
	squaredFigChan := make(chan Cube, 2)

	go computeArea(cubes, squaredFigChan)

	sendCubes(squaredFigChan)
}

func computeArea(cubes []Cube, squaredFigChan chan<- Cube) {
	for _, cube := range cubes {
		cube.Area = cube.Length * cube.Width
		squaredFigChan <- cube
	}
	close(squaredFigChan)
}

func sendCubes(squaredFigChan <-chan Cube) {
	for cube := range squaredFigChan {
		fmt.Println(cube)
	}
}
