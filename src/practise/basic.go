package practise

import (
	"fmt"
	"math/rand"
	"time"
)

func Run() {
	// nil 只能和对象，数组，指针做比较

	//testArray()
	//testPointer()
	//testStructure()

	//testSlice()
	//testRange()

	//testConvert()

	//testInterface()

	testConcurrent()
}



func printTitle(title string) {
	fmt.Printf("-------%s------\n", title)
}
func testArray() {
	printTitle("test array")
	//数组 类型相同，长度固定
	//指定长度的数组只能传值给形参指定数组长度的函数
	var a = [5]int{1, 2, 3, 4, 5}
	fmt.Println(a)

	var b = [...]float32{1, 2.0, 3.33, 4.69}
	fmt.Println(b)

	var c [10]int
	for i := 0; i < len(c); i++ {
		c[i] = i
	}
	fmt.Println(c)
}

func testPointer()  {
	printTitle("array pointer")
	//指针
	var a = 10
	fmt.Printf("变量a的地址：%x\n", &a)

	var ptr *int
	if ptr == nil {
		fmt.Println("未赋值的指针是nil", ptr)
	} else {
		fmt.Println("未赋值的指针不是nil", ptr)
	}
	fmt.Printf("未赋值的指针的值：%x\n", ptr)

	ptr = &a
	fmt.Printf("赋值 地址：%x, 值：%d\n", ptr, *ptr)

	var b int
	if &b == nil {
		fmt.Println("未赋值的变量的地址是nil", &b)
	} else {
		fmt.Println("未赋值的变量的地址不是nil", &b)
	}

	var x = 100
	var y = 200
	fmt.Println("交换前的x,y: ", x, y)
	swap(&x, &y)
	fmt.Println("交换后的x,y: ", x, y)
}

func swap(x *int, y *int){
	*x, *y = *y, *x
}


func testStructure(){
	printTitle("array struct")
	var c Config1
	var c2 = Config1{}
	c2.name = "-"
	var c3 = Config1{"AAA", 12, 0, false}
	fmt.Println("未实例化的结构体 ", c)
	fmt.Println("空的实例化结构体 ", c2)
	fmt.Println("实例化的结构体 ", c3)

}
type Config1 struct {
	name string
	id uint64
	gender uint8
	visitor bool
}

func testSlice() {
	printTitle("test slice")

	var a = make([]int, 5, 10)
	fmt.Println(a, len(a), cap(a))
	a = append(a, 1,2,3)
	fmt.Println(a)
	fmt.Println(a[4:])

	var b = make([]int, 5)
	copy(b, a)
	b[1] = 5
	fmt.Println(b)
	fmt.Println(a)
}

func testRange(){
	printTitle("test range")

	//在for循环中迭代 数组，切片，集合，通道的元素
	//对于数组，切片， 返回的是index, value
	//对于集合，返回的是key,value

	var a = [5]int{1,2,3,4,5}
	for i,v := range a {
		fmt.Println("数组迭代:", i, v)
	}

	var b = make(map[string]int)
	b["a"] = 65
	b["b"] = 'B'
	b["c"] = 'C'
	b["d"] = 'D'
	//map 迭代 顺序随机
	for k, v := range b {
		fmt.Println("map迭代：", k, v)
	}

	var c = map[string]string{"A":"1", "B":"2", "C":"3"}
	for k,v := range c {
		fmt.Println("map迭代2：", k, v)
	}
}


func testConvert(){
	printTitle("test convert")

	var sum = 17
	var count  = 5
	var mean float32

	mean = float32(sum) / float32(count)
	fmt.Println("mean的值：", mean)

	fmt.Printf("%f: \n", 12.5)
	fmt.Printf("int=>float %f: \n", float32(12))
	//只有 低精度 可以向 高精度 转换
	//int => float32 => float64
	//fmt.Printf("float=>int %d: \n", int(12.545))
}


type Phone interface {
	call()
	getName() string
	getPrice() int
}

type NokiaPhone struct {
	name string
}

type IPhone struct {
	os string
}
func (iPhone IPhone) call() {
	iPhone.os = "iOS"
	fmt.Println("this is iPhone, os is:", iPhone.os)
}
func (iPhone IPhone) getName() string {
	return "iPhone 8"
}
func (iPhone IPhone) getPrice() int {
	return 6500
}

func (nokiaPhone NokiaPhone) call() {
	nokiaPhone.name = "nokiaPhone"
	fmt.Println("this is", nokiaPhone.name)
}
func (nokiaPhone NokiaPhone)  getName() string {
	return nokiaPhone.name
}
func (nokiaPhone NokiaPhone)  getPrice() int {
	return 3200
}

func testInterface()  {
	printTitle("test interface")

	var phone Phone
	phone = new(NokiaPhone)
	phone.call()
	fmt.Println(phone.getName(),"的价格", phone.getPrice())

	phone = new(IPhone)
	phone.call()
	fmt.Println(phone.getName(),"的价格", phone.getPrice())

}

/**
 * 测试并发
 */
func testConcurrent()  {
	printTitle("test concurrent")

	//被go修饰的方法只能放在顶部执行
	//go say("1 World")
	//go say("2 Hi")
	//say("3 Hello")

	var t0 = time.Now();
	const size = 300000000
	var a = make([]uint64, size)
	for i:=0;i<size;i++ {
		a[i] = rand.Uint64() / 100
	}
	fmt.Println("赋值耗时：", time.Since(t0))

	//单线程
	var t1 = time.Now()
	s1 := sum2(a)
	var d1 = time.Since(t1)
	fmt.Println("耗时d1", d1, s1) //这个时间其实并不准确，当size过大时

	//多线程 之后多线程时才能用管道(不确定)
	var t2 = time.Now()
	//只能这样格式化化，使用golang诞生的时间
	var c2 = make(chan uint64)
	var golen = 5
	var gap = size / golen
	for i := 0; i < golen; i++ {
		go sum(a[gap*i: gap*(i+1)], c2)
	}
	var s2 uint64
	for i:=0;i<golen;i++ {
		s2 += <-c2
	}
	//var d2 = time.Since(t2)
	var t3 = time.Now()
	fmt.Println("耗时d2", t3.Sub(t2), s2) //这个时间其实并不准确，当size过大时，貌似得×100

}

func say(s string) {
	for i:=0; i<5;i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []uint64, c chan uint64){
	var sum uint64 = 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}
func sum2(s []uint64) uint64 {
	var sum uint64 = 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func testString() {

}