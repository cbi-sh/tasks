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

	checkLink  string
	checkValue string

	actionLink string
}

func daySeconds(t time.Time) uint32 {
	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return uint32(t.Sub(midnight).Seconds())
}

func main() {

	done := make(chan bool)
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Println("Hello !!")
				perform(Task{
					name:       "Greenhouse light off",
					start:      0 * 3600,
					stop:       5*3600 + 59*60 + 59,
					checkLink:  "http://192.168.0.21/light",
					checkValue: "0",
					actionLink: "http://192.168.0.21/light/off",
				})
				fmt.Println("Hello !@")
				perform(Task{
					name:       "Greenhouse light on",
					start:      6 * 3600,
					stop:       23*3600 + 59*60 + 59,
					checkLink:  "http://192.168.0.21/light",
					checkValue: "1",
					actionLink: "http://192.168.0.21/light/on",
				})
			}
		}
	}()

	time.Sleep(1000000 * time.Second)
	done <- true
}

func perform(task Task) {

	now := daySeconds(time.Now())

	if task.start <= now && now < task.stop {
		fmt.Printf("checking task: %s\n", task.name)

		resp, err := http.Get(task.checkLink)
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

		if string(body) != task.checkValue {
			fmt.Println("value don't match, try to make action")

			resp, err := http.Get(task.actionLink)
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
