package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/mkevac/debugcharts"
)

var (
	clientNums    = 3
	messageNums   = 100
	dataStreaming []chan data

	token int64
	step  int64

	maxSleepInterval int64 = 5
	maxGap           int64 = 10

	wg sync.WaitGroup

	debug            = false
	enabledRateLimit = true
	rate             = time.Second / 100
)

type data struct {
	kind    string
	prepare int64
	commit  int64
}

func init() {
	dataStreaming = make([]chan data, clientNums)
	for i := 0; i < clientNums; i++ {
		dataStreaming[i] = make(chan data, messageNums)
	}
	if debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}

/* 1. 请实现 collectAndSort 函数：假设数据源产生的数据是无休无止的，请以最快的方法输出排序后的结果
 * 2. 考虑流控，控制内存的使用
 */
func main() {
	wg.Add(clientNums*2 + 1)
	for i := 0; i < clientNums; i++ {
		go func(index int) {
			defer wg.Done()
			generateDatas(index)
		}(i)
		go func(index int) {
			defer wg.Done()
			generateDatas(index)
		}(i)
	}

	go func() {
		defer wg.Done()
		collectAndSort()
	}()
	wg.Wait()
}

func collectAndSort() {
	dataChan := make(chan data, clientNums*2*messageNums)
	var throttle *time.Ticker
	go func() {
		if enabledRateLimit {
			defer throttle.Stop()
			throttle = time.NewTicker(rate)
		}

		for {
			for _, c := range dataStreaming {
				if enabledRateLimit {
					<-throttle.C
				}
				data := <-c
				// 过滤出Commit的数据
				if data.kind == "commit" {
					dataChan <- data
				}
			}
		}
	}()

	var dataSet []data

	for {
		select {
		case data := <-dataChan:
			// 1、收集数据
			dataSet = append(dataSet, data)

			if len(dataSet) == clientNums*2*messageNums {
				// 2、数据排序
				MyMergeSort(dataSet)
				// 3、数据打印
				for _, d := range dataSet {
					fmt.Println(d)
				}
			}
		}
	}
}

func generateDatas(index int) {
	for i := 0; i < messageNums; i++ {
		prepare := incrementToken()
		sleep(maxSleepInterval)

		dataStreaming[index] <- data{
			kind:    "prepare",
			prepare: prepare,
		}
		sleep(maxSleepInterval)

		commit := incrementToken()
		sleep(maxSleepInterval)

		dataStreaming[index] <- data{
			kind:    "commit",
			prepare: prepare,
			commit:  commit,
		}
		sleep(10 * maxSleepInterval)
	}
}

func incrementToken() int64 {
	token := atomic.AddInt64(&token, rand.Int63()%maxGap+1)
	return token
}

func sleep(factor int64) {
	interval := atomic.AddInt64(&step, 3)%factor + 1
	waitTime := time.Duration(rand.Int63() % interval)
	time.Sleep(waitTime * time.Millisecond)
}
