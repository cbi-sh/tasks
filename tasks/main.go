package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Task struct {
	name string

	start uint32
	stop  uint32

	check_link  string
	check_value string

	action_link string
}

func daySeconds(t time.Time) uint32 {
	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return uint32(t.Sub(midnight).Seconds())
}

func main() {

	done := make(chan bool)
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Println("Hello !!")
				check(Task{
					name:        "Greenhouse light on",
					start:       6 * 3600,
					stop:        23*3600 + 59*60 + 59,
					check_link:  "http://192.168.0.21/light",
					check_value: "1",
					action_link: "http://192.168.0.21/light/on",
				})
			}
		}
	}()

	time.Sleep(10000 * time.Second)
	done <- true
}

func check(task Task) {

	fmt.Printf("%+v\n", task)
	fmt.Printf("%+v\n", daySeconds(time.Now()))

	now := daySeconds(time.Now())

	if task.start <= now && now < task.stop {
		fmt.Printf("checking task: %s\n", task.name)

		resp, err := http.Get(task.check_link)
		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
			return
		} else {
			resp.Body.Close()
		}

		value := string(body)

		fmt.Println(value)

		if value != task.check_value {

			fmt.Println("value don't match, try to make action")

			resp, err := http.Get(task.action_link)
			if err != nil {
				fmt.Println(err)
				return
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Println(err)
				return
			} else {
				resp.Body.Close()
			}

			value := string(body)

			fmt.Println("answer from device: ", value)
		}
	}
}
