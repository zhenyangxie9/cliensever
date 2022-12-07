package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"uk.ac.bris.cs/gameoflife/Broker/newstubs"
	"uk.ac.bris.cs/gameoflife/util"
)

var CurWorld [][]uint8
var UpdateWorld [][]uint8
var CurTurn int
var mu sync.Mutex

type Worker struct{}

func calculateNextState(startY int, endY int, startX int, endX int, world [][]uint8) [][]uint8 {
	ImageHeight := endY - startY
	ImageWidth := endX - startX
	newWorld := NewWorld(ImageHeight, ImageWidth)
	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			sum := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					if world[(y+startY+ImageHeight+i)%ImageHeight][(x+ImageWidth+j)%ImageWidth] == 255 {
						//if world[(y+startY+i)%startY][(x+width+j)%width] == 255 {
						sum += 1
					}
				}
			}
			if world[y+startY][x] == 255 {
				sum--
			}
			if world[y+startY][x] == 255 {
				if sum < 2 {
					newWorld[y][x] = 0
					//c.events <- CellFlipped{turn, util.Cell{y + startY, x}}
				} else if sum == 2 || sum == 3 {
					newWorld[y][x] = 255
				} else if sum > 3 {
					newWorld[y][x] = 0
					//c.events <- CellFlipped{turn, util.Cell{y + startY, x}}
				}
			} else if world[y+startY][x] == 0 {
				if sum == 3 {
					newWorld[y][x] = 255
					//c.events <- CellFlipped{turn, util.Cell{y + startY, x}}
				} else {
					newWorld[y][x] = 0
				}
			}
		}
	}
	return newWorld
}
func NewWorld(height int, width int) [][]uint8 {
	world := make([][]uint8, height)
	for i := range world {
		world[i] = make([]uint8, width)
	}
	return world
}
func calculateAliveCells(ImageWidth, ImageHeight int) []util.Cell {
	aliveCell := make([]util.Cell, 0)
	for x := 0; x < ImageWidth; x++ {
		for y := 0; y < ImageHeight; y++ {
			if CurWorld[x][y] == 255 {
				aliveCell = append(aliveCell, util.Cell{X: y, Y: x})
			}
		}
	}
	return aliveCell
}

func (w *Worker) NextState(req newstubs.BrokerWRequest, res *newstubs.BrokerWResponse) (err error) {
	startY := res.StartY
	endY := res.EndY
	startX := 0
	endX := res.ImageWidth
	CurWorld = NewWorld(res.ImageHeight, res.ImageWidth)
	workers := res.Workers
	for i := 0; i < workers; i++ {
		//if i == workers-1 {
		//	go calculateNextState(startY, endY, startX, endX, CurWorld)
		//} else {
		//	go calculateNextState(startY, endY, startX, endX, CurWorld)
		//}
		go calculateNextState(startY, endY, startX, endX, CurWorld)
	}

	//res.World=
	return
}

//func worker(workers int, WorkerClient []*rpc.Client, ImageHeight int, ImageWidth int) {
//	for i := 0; i < workers; i++ {
//		startY := i * ImageHeight / workers
//		var endY int
//		if i == workers-1 {
//			endY = ImageHeight
//		} else {
//			endY = (i + 1) * ImageHeight
//		}
//		req := stubs.BrokerWRequest{StartY: startY, EndY: endY, ImageWidth: ImageWidth, ImageHeight: ImageHeight, World: CurWorld}
//		var res []*stubs.BrokerWResponse
//		res = append(res, new(stubs.BrokerWResponse))
//
//		WorkerClient[].Go()
//	}
//}

func main() {
	pAddr := flag.String("port", "8050", "Port to listen on")
	flag.Parse()
	//rand.Seed(time.Now().UnixNano())
	rpc.Register(&Worker{})

	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("listening on %s", listener.Addr().String())
	defer listener.Close()
	rpc.Accept(listener)
}
