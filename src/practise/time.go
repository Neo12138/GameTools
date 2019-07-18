package practise

import (
	"fmt"
	"time"
)

func RunTime() {
	testSince()
}

func testSince() {
	printTitle("test time.since")

	t1 := time.Date(2019,4,12,12,0,0,0, time.Local)
	fmt.Println(t1)
	t2 := time.Date(2019,4,12,12,1,1,44444, time.Local)
	fmt.Println(t2)
	fmt.Println(t2.Sub(t1))

	t3 := time.Now()
	fmt.Println(t3)
	time.Sleep(time.Second * 20)
	t4 := time.Now()
	fmt.Println(t4.Sub(t3))
	fmt.Println(time.Since(t3))
}


