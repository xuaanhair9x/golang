package main

import(
	"fmt"
	"time"
	"sync"
)

type Fork struct {
	id int
	l  chan bool  
}

type Philosopher struct {
	id                  int
	leftFork, rightFork *Fork // 2 chiếc dĩa bên cạnh
}

var ForkCond *sync.Cond

func (p *Philosopher) pickForks() {
	ForkCond.L.Lock() // lock mutex
	for {
	   select {
	   case <-p.leftFork.l: 
		  fmt.Println(p.id, " picked left fork ", p.leftFork.id)
		  <-p.rightFork.l
		  fmt.Println(p.id, " picked right fork ", p.rightFork.id)
		  ForkCond.L.Unlock() // unlock mutex
		  return
	   default:
		  ForkCond.Wait() //wait cho tới khi được signaled
	   }
	}
 }

func (p *Philosopher) putDownForks() {
	p.leftFork.l <- true // chiếc đũa bên trái về trạng thái ready trong channel
	fmt.Println(p.id, " put down left fork ", p.rightFork.id)
	p.rightFork.l <- true // chiếc đũa bên phải về trạng thái ready trong channel
	fmt.Println(p.id, " put down right fork ", p.leftFork.id)

	ForkCond.Signal()
}
const (
	HUNGER = 3
	COUNT  = 5
)

var (
	wg       sync.WaitGroup
	forks    [COUNT]*Fork
)

func main() {
	var mu sync.Mutex
	ForkCond = sync.NewCond(&mu)

	wg.Add(COUNT) 
	for i := 0; i < COUNT; i++ { // khởi tạo forks
		forks[i] = &Fork{id: i, l: make(chan bool, 1)}
		forks[i].l <- true
	}

	philosophers := make([]*Philosopher, COUNT)
	for i := 0; i < COUNT; i++ {
		philosophers[i] = &Philosopher{
			id: i, leftFork: forks[i], rightFork: forks[(i+1)%COUNT]}
		go philosophers[i].dine() // bắt đầu bữa tối
	}

	wg.Wait() // chờ philosophers hoàn thành bữa tối
	fmt.Println("Table is empty")
}

func (p *Philosopher) dine() {

	for i := 0; i < HUNGER; i++ {
		say("Thinking", p.id)
		doingStuff()

		say("Hungry", p.id)
		p.pickForks()

		say("Eating", p.id)
		doingStuff()
		p.putDownForks()
	}
	say("Satisfied", p.id)
	wg.Done() // hoàn thành bữa tối
	say("Leaving the table", p.id)
}
func doingStuff() {
	time.Sleep(time.Second / 10)
}

func say(action string, id int) {
	fmt.Printf("#%d is %s\n", id, action)
}