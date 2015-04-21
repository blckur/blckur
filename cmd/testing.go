package cmd

import (
	"fmt"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/utils"
	"strconv"
	"time"
)

// Starts test publisher
func Pub() {
	i := 0

	for {
		conn := cache.Get()
		err := conn.Publish("test", fmt.Sprintf("msg-%d", i), nil)
		if err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(500 * time.Millisecond)
		i += 1
	}
}

// Starts test subscriber
func Sub() {
	lst := cache.Subscribe("test")

	for msg := range lst.Listen() {
		fmt.Println(msg.Type)
	}
}

// Starts test publish and subscribe
func PubSub() {
	lst := cache.Subscribe("test")

	go func() {
		for msg := range lst.Listen() {
			fmt.Println(msg)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	conn := cache.Get()
	err := conn.Publish("test", "msg", nil)
	if err != nil {
		panic(err)
	}
}

// Starts publish performance test
func StressPub() {
	start := time.Now()

	wait := utils.WaitCancel{}
	wait.Add(5)

	for j := 0; j < 5; j++ {
		go func(j int) {
			conn := cache.Get()

			for i := 20000 * j; i < 20000*(j+1); i++ {
				x := strconv.Itoa(i)
				err := conn.Publish(x, x, nil)
				if err != nil {
					fmt.Printf("err: %d\n", i)
				}
			}

			wait.Done()
		}(j)
	}

	wait.Wait()

	fmt.Printf("publish time: %s\n", time.Since(start).String())
}

// Starts subscribe performance test
func StressSub() {
	start := time.Now()

	wait := utils.WaitCancel{}
	wait.Add(100000)

	for i := 0; i < 100000; i++ {
		go func(i int) {
			channel := fmt.Sprintf("%d", i)
			lst := cache.Subscribe(channel)

			go func() {
				for msg := range lst.Listen() {
					if i%10000 == 0 || i == 99999 {
						fmt.Printf("%s: %s\n", channel, msg.Type)
					}
				}
			}()

			wait.Done()
		}(i)
	}

	wait.Wait()

	fmt.Printf("subscribe time: %s\n", time.Since(start).String())

	time.Sleep(6 * time.Hour)
}
