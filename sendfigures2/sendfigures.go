package sendfigures2

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
	cubesWithAreaChan := make(chan Cube, 2)
	cubesWithVolumeChan := make(chan Cube, 2)

	go computeArea(cubes, cubesWithAreaChan)
	go computeVolume(cubesWithAreaChan, cubesWithVolumeChan)

	sendCubes(cubesWithVolumeChan)
}

func computeArea(cubes []Cube, cubesWithAreaChan chan<- Cube) {
	for _, cube := range cubes {
		cube.Area = cube.Length * cube.Width
		cubesWithAreaChan <- cube
	}
	close(cubesWithAreaChan)
}

func computeVolume(cubesWithAreaChan <-chan Cube, cubesWithVolumeChan chan<- Cube) {
	for cube := range cubesWithAreaChan {
		cube.Volume = cube.Area * cube.Height
		cubesWithVolumeChan <- cube
	}
	close(cubesWithVolumeChan)
}

func sendCubes(cubesWithVolumeChan <-chan Cube) {
	for cube := range cubesWithVolumeChan {
		fmt.Println(cube)
	}
}
