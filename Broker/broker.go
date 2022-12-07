package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"uk.ac.bris.cs/gameoflife/Broker/newstubs"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

var GloWorld [][]uint8
var GloTurn int
var GloThreads int
var mu sync.Mutex

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

func calculateAliveCells(ImageWidth, ImageHeight int, world [][]uint8) []util.Cell {
	aliveCell := make([]util.Cell, 0)
	for x := 0; x < ImageWidth; x++ {
		for y := 0; y < ImageHeight; y++ {
			if world[x][y] == 255 {
				aliveCell = append(aliveCell, util.Cell{X: y, Y: x})
			}
		}
	}
	return aliveCell
}

func NewWorld(height int, width int) [][]uint8 {
	newWorld := make([][]uint8, height)
	for i := range newWorld {
		newWorld[i] = make([]uint8, width)
	}
	return newWorld
}

//func World(c distributorChannels, p gol.Params) [][]uint8 {
//	width := p.ImageWidth
//	height := p.ImageHeight
//	world := NewWorld(p.ImageHeight, p.ImageWidth)
//	for h := 0; h < height; h++ {
//		for w := 0; w < width; w++ {
//			world[h][w] = <-c.ioInput
//		}
//	}
//	return world
//}

func worker(startY, endY, startX, endX int, out chan<- [][]uint8, world [][]uint8) {
	newWorld := calculateNextState(startY, endY, startX, endX, world)
	out <- newWorld
}

//func Advance(Width, Height, Threads int) {
//	client, _ := rpc.Dial("tcp", "127.0.0.1:8050")
//	defer client.Close()
//	Len := Width / Threads
//	for i := 0; i < Threads; i++ {
//		startY := i * Len
//		endY := (i + 1) * Len
//		if i == Threads-1 {
//			endY = Height
//		}
//		req := stubs.BrokerWRequest{Turns: GloTurn, ImageWidth: Width, ImageHeight: Height, World: GloWorld}
//		res := stubs.BrokerWResponse{}
//		client.Call(stubs.NetState, req, res)
//		GloWorld = calculateNextState(0, res.ImageHeight, 0, res.ImageWidth, GloWorld)
//	}
//}

func workers(newWorld [][]uint8, Threads, ImageHeight, ImageWidth int) [][]uint8 {

	if Threads == 1 {
		newWorld = calculateNextState(0, ImageHeight, 0, ImageWidth, newWorld)
	} // else {
	//	out := make([]chan [][]uint8, Threads)
	//	for i := range out {
	//		out[i] = make(chan [][]uint8)
	//	}
	//	for i := 0; i < Threads; i++ {
	//		if i == Threads-1 {
	//			go worker(i*ImageHeight/Threads, ImageHeight, 0, ImageWidth, out[i], newWorld)
	//		} else {
	//			go worker(i*ImageHeight/Threads, (i+1)*ImageHeight/Threads, 0, ImageWidth, out[i], newWorld)
	//		}
	//		//go worker(i*p.ImageHeight/p.Threads, (i+1)*p.ImageHeight/p.Threads, out[i], newWorld, p)
	//	}
	//
	//	newWorld = nil
	//	//newWorld = NewWorld(0, 0)
	//	for i := 0; i < Threads; i++ {
	//		part := <-out[i]
	//		//fmt.Println(part)
	//		newWorld = append(newWorld, part...)
	//	}
	//}
	return newWorld

}

func Update(workers int, Width int, Height int, Clients *rpc.Client) {
	var out []chan *rpc.Call
	Len := Height / workers
	for i := 0; i < workers; i++ {
		startY := i * Len
		endY := (i + 1) * Len
		if i == workers-1 {
			endY = Height
		}
		req := newstubs.BrokerWRequest{StartY: startY, EndY: endY, ImageWidth: Width, ImageHeight: Height, World: GloWorld, Workers: workers}
		var ress []*newstubs.BrokerWResponse
		ress = append(ress, new(newstubs.BrokerWResponse))
		out = append(out, make(chan *rpc.Call, 1))
		//Clients[i].Go(stubs.NextState, req, ress[i], out[i])
		Clients.Go(newstubs.NextState, req, ress[i], out[i])
	}
	for i := 0; i < workers; i++ {
		<-out[i]
	}

}

type GameOfLife struct{}

func (s *GameOfLife) ProcessGol(req stubs.Request, res *stubs.Response) (err error) {
	client, _ := rpc.Dial("tcp", "127.0.0.1:8050")
	defer client.Close()
	GloWorld = req.World
	GloTurn = 0
	GloThreads = req.Threads
	for GloTurn < req.Turns {
		GloTurn++
		mu.Lock()
		//GloWorld = calculateNextState(0, req.ImageHeight, 0, req.ImageWidth, GloWorld)
		//GloWorld =
		Update(GloThreads, req.ImageWidth, req.ImageHeight, client)
		mu.Unlock()

	}
	res.World = GloWorld
	res.Turns = GloTurn
	return
}

func (s *GameOfLife) AliveCell(req stubs.Request, res *stubs.Response) (err error) {
	mu.Lock()
	res.World = GloWorld
	res.Turns = GloTurn
	mu.Unlock()
	return
}

// keypress s,p,q,k in stage 2
func (s *GameOfLife) ShutDown(req stubs.Request, res *stubs.Response) (err error) {
	os.Exit(0)
	return
}

func (s *GameOfLife) PauseServer(req stubs.Request, res *stubs.Response) (err error) {
	res.Turns = GloTurn
	mu.Lock()
	return
}
func (s *GameOfLife) ResetServer(req stubs.Request, res *stubs.Response) (err error) {
	mu.Unlock()
	return
}
func (s *GameOfLife) CurrWorld(req stubs.Request, res *stubs.Response) (err error) {
	res.World = GloWorld
	res.Turns = GloTurn
	return
}
func (s *GameOfLife) CloseDis(req stubs.Request, res *stubs.Response) (err error) {
	closeDis := make(chan bool)
	closeDis <- true
	return
}

//Broker

func main() {
	pAddr := flag.String("port", "8040", "Port to listen on")
	flag.Parse()
	//rand.Seed(time.Now().UnixNano())
	rpc.Register(&GameOfLife{})

	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("listening on %s", listener.Addr().String())
	defer listener.Close()

	rpc.Accept(listener)
}
