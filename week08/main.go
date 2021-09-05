package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"runtime"
	"strconv"
	"time"
)

func main()  {
	redis := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   1000,
		IdleTimeout: 60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
	defer redis.Close()

	reConn := redis.Get()
	defer reConn.Close()

	// init 10000 size "1"
	var value string = "1"
	for j := 0; j < 10000; {
		value += "1"
		j++
	}

	var avg uint32
	for i := 0; i < 30; i++ {
		var tempValue = ""
		beforeMem := MemStat()

		// init 1w, 2w, 3w ..... 30w
		for k := -1; k < i; k++ {
			tempValue += value
		}

		//fmt.Println("插入到redis数据：key= key", strconv.Itoa(i), ",value.size= ", strconv.Itoa(i+1)+"W", " ,real size= ", len(tempValue))

		_, err := reConn.Do("set", "key"+strconv.Itoa(i+1), tempValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		tempValue = ""

		// get memery info
		time.Sleep(1 * time.Second)
		afterMem := MemStat()

		var used uint32 = 0
		if afterMem.Used < beforeMem.Used {
			used = 0
		} else {
			used = afterMem.Used - beforeMem.Used
		}

		fmt.Println(
			"redis数据：插入前key为："+strconv.Itoa(i+1),
			" 内存使用情况: ", beforeMem.Used,
			"b, 插入后key", strconv.Itoa(i+1),
			" 内存使用情况: ", afterMem.Used,
			"b, 插入占用内存: ", used, "b\n")

		avg += used
	}

	avg = avg / 50 / 1000 // k
	fmt.Println("\n50个key平均占用内存 ", avg, "k")

	// clear all redis keys
	for i := 0; i < 50; i++ {
		_, err := reConn.Do("del", "key"+strconv.Itoa(i+1))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
type MemStatus struct {
	All  uint32 `json:"all"`
	Used uint32 `json:"used"`
	Free uint32 `json:"free"`
	Self uint64 `json:"self"`
}

func MemStat() MemStatus {
	// 自身占用
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	mem := MemStatus{}
	mem.Self = memStat.Alloc

	// 内存情况 win10无法获取

	return mem
}
