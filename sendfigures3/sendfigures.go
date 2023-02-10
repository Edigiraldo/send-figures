package sendfigures3

import (
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Cube struct {
	Length int
	Width  int
	Height int
	Area   int
	Volume int
}

var (
	statusOk     = 0
	statusFailed = 1
)

func Process() {
	programStatus := statusOk
	cubes := []Cube{{1, 2, 5, 0, 0}, {4, 6, 7, 0, 0}, {8, 10, 2, 0, 0}, {7, 2, 7, 0, 0}, {6, 1, 8, 0, 0}}

	cubesWithAreaChan := make(chan Cube, 2)
	cubesWithVolumeChan := make(chan Cube, 2)
	errChan := make(chan error)

	errWaitGroup := sync.WaitGroup{}
	errWaitGroup.Add(1)

	go func() {
		for err := range errChan {
			programStatus = statusFailed
			fmt.Printf("the process failed: %v\n", err)
		}
		errWaitGroup.Done()
	}()

	go func() {
		if err := computeArea(cubes, cubesWithAreaChan); err != nil {
			errChan <- err
		}
	}()
	go func() {
		if err := computeVolume(cubesWithAreaChan, cubesWithVolumeChan); err != nil {
			errChan <- err
		}
	}()

	sendCubes(cubesWithVolumeChan, errChan)
	fmt.Println("the program finished with status: ", programStatus)
}

func computeArea(cubes []Cube, cubesWithAreaChan chan<- Cube) (err error) {
	errGroup := errgroup.Group{}
	defer close(cubesWithAreaChan)

	for _, cube := range cubes {
		cubeClosure := cube
		errGroup.Go(func() error {
			if cubeClosure.Length < 0 || cubeClosure.Width < 0 {
				return fmt.Errorf("the length neither width in cube %v cannot be negative -> length(%d), width(%d)",
					cubeClosure, cubeClosure.Length, cubeClosure.Width)
			}

			cubeClosure.Area = cubeClosure.Length * cubeClosure.Width
			cubesWithAreaChan <- cubeClosure
			return nil
		})
	}

	return errGroup.Wait()
}

func computeVolume(cubesWithAreaChan <-chan Cube, cubesWithVolumeChan chan<- Cube) (err error) {
	errGroup := errgroup.Group{}
	defer close(cubesWithVolumeChan)

	for cube := range cubesWithAreaChan {
		cubeClosure := cube
		errGroup.Go(func() error {
			if cubeClosure.Height < 0 {
				return fmt.Errorf("the height in cube %v cannot be negative -> height(%d)",
					cubeClosure, cubeClosure.Height)
			}

			cubeClosure.Volume = cubeClosure.Area * cubeClosure.Height
			cubesWithVolumeChan <- cubeClosure
			return nil
		})
	}

	return errGroup.Wait()
}

func sendCubes(cubesWithVolumeChan <-chan Cube, errChan chan<- error) {
	for cube := range cubesWithVolumeChan {
		var err error = nil
		if cube.Volume > 150 {
			err = fmt.Errorf("was not possible to send cube(%v) with volume=%d greater than 150", cube, cube.Volume)
			errChan <- err
		}
		if err == nil {
			fmt.Println(cube, "-------", "_/")
		}
	}
}
