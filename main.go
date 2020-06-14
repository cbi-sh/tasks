package main

import (
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

func daySeconds(t time.Time) uint32 {
	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return uint32(t.Sub(midnight).Seconds())
}

func (task *Task) perform() {

	log.Info("perform task: ", task.name)

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

				log.Info("start checking task")

				(&Task{
					name:       "Greenhouse light off",
					start:      0 * 3600,
					stop:       5*3600 + 59*60 + 59,
					checkLink:  "http://192.168.0.21/light",
					checkValue: "0",
					actionLink: "http://192.168.0.21/light/off",
				}).perform()

				(&Task{
					name:       "Greenhouse light on",
					start:      6 * 3600,
					stop:       23*3600 + 59*60 + 59,
					checkLink:  "http://192.168.0.21/light",
					checkValue: "1",
					actionLink: "http://192.168.0.21/light/on",
				}).perform()
			}
		}
	}()

	time.Sleep(1000000 * time.Second)
	done <- true
}
