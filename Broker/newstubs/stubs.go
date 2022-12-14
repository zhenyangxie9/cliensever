package newstubs

var Gameoflife = "GameOfLife.ProcessGol"
var AliveCell = "GameOfLife.AliveCell"
var CurrentState = "GameOfLife.CurrWorld"
var CloseDis = "GameOfLife.CloseDis"
var ShutDown = "GameOfLife.ShutDown"
var Pause = "GameOfLife.PauseServer"
var Reset = "GameOfLife.ResetServer"

var NextState = "Worker.NextState"

type Response struct {
	World       [][]uint8
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
	Pause       chan bool
}

type Request struct {
	World       [][]uint8
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

type BrokerWResponse struct {
	StartY int
	EndY   int
	//Turns       int
	ImageWidth  int
	ImageHeight int
	World       [][]uint8
	Workers     int
}

type BrokerWRequest struct {
	StartY int
	EndY   int
	//Turns       int
	ImageWidth  int
	ImageHeight int
	World       [][]uint8
	Workers     int
}
