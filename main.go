package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Task struct {
	name string

	start uint32
	stop  uint32

	checkLink  string
	checkValue string

	actionLink string
}

func (task Task) Perform() {

	now := daySeconds(time.Now())
	if !(task.start <= now && now < task.stop) {
		log.Info("skip task by time: ", task.name)
		return
	}

	resp, err := http.Get(task.checkLink)
	if err != nil {
		log.Error("check link error: ", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read body error: ", err)
		return
	}
	resp.Body.Close()

	if string(body) == task.checkValue {
		log.Info("value ok, ", task.name)
		return
	}

	log.Info("value don't match, try to make action")
	resp, err = http.Get(task.actionLink)
	if err != nil {
		log.Error("check link error: ", err)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read body error: ", err)
		return
	}
	resp.Body.Close()

	log.Info("answer from device: ", string(body))
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
		tasks := []Task{
			{
				name:       "Hold light off",
				start:      0 * 3600,
				stop:       8*3600 + 59*60 + 59,
				checkLink:  "http://192.168.0.21/light",
				checkValue: "0",
				actionLink: "http://192.168.0.21/light/off",
			},
			{
				name:       "Hold light on",
				start:      9 * 3600,
				stop:       23*3600 + 59*60 + 59,
				checkLink:  "http://192.168.0.21/light",
				checkValue: "1",
				actionLink: "http://192.168.0.21/light/on",
			},
		}

		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				for _, task := range tasks {
					task.Perform()
				}
			}
		}
	}()

	fmt.Scanln()
	done <- true
}
