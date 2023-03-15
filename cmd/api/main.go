package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	rt "github.com/papaburgs/almagest/pkg/redistools"
)

var arc *rt.AlmagestRedisClient

type statusItem struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Health  string `json:"health"`
}

type statusList struct {
	Status   string       `json:"status"`
	Services []statusItem `json:"services"`
}

type discordBody struct {
	Channel string
	Content string
}

//go:embed gitc.txt
var gitCommit string

func main() {

	arc = rt.New()
	go arc.PublishWatchdog("api")
	log.SetLevel(log.DebugLevel)
	// this will be the base mux
	m := http.NewServeMux()

	m.HandleFunc("/api/almagest/discord/dispatch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		switch r.Method {
		case http.MethodGet:
			log.Debug("function called", "endpoint", "discord dispatch", "method", "GET")
			res := []byte(`{"status":"success", "message":"In order to send a message to discord, send post message",
			        "schema":{"channel": "botspot", "message": "post me", "server": "32ohsix"}}`)
			w.Write(res)

		case http.MethodPost:
			log.Debug("function called", "endpoint", "discord dispatch", "method", "POST")
			var b discordBody
			content, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Error("error reading content")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if len(content) == 0 {
				http.Error(w, `{"status": "fail: no content"}`, http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(content, &b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dsm := rt.PSMessage{
				Service:   "discord",
				MessageID: uuid.New().String(),
				Channel:   b.Channel,
				Content:   b.Content,
			}
			err = arc.Publish(dsm)
			if err != nil {
				log.Error("error posting to redis", "error", err)
			}
			res := `{"status":"success","message":"message dispatched to discord"}`
			w.Write([]byte(res))
			log.Info("published to redis")

		default:
			http.Error(w, "Unknown method", http.StatusNotImplemented)
			return
		}
	})

	m.HandleFunc("/api/almagest/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		c := arc.Subscribe()
		myMessageID := uuid.NewString()
		dl := statusList{Status: "success", Services: []statusItem{}}

		// send status request message
		dsm := rt.PSMessage{
			Service:   "healthcheck",
			MessageID: myMessageID,
		}
		arc.Publish(dsm)

		//now respond to my own message
		log.Debug("replying to health check request")
		arc.PostStatus("api", strings.TrimSpace(gitCommit), myMessageID)

		// waiting 15 seconds or the responses
		timer := time.NewTimer(15 * time.Second)
		timerDone := true

		for timerDone {
			log.Debug("starting select")
			select {
			case <-timer.C:
				log.Debug("timer done, returning what I have")
				timerDone = false
				break
			case msg := <-c:
				log.Debug("picked up a message")

				psm, class, err := rt.ClassifyMessage(msg)
				if err != nil {
					log.Error("Could not decode message ", "payload", msg.Payload, "error", err)
					continue
				}
				if class == rt.HealthCheckResponse && psm.ResponseTo == myMessageID {
					content := strings.Split(psm.Content, "|")

					if len(content) < 3 {
						log.Error("Content was not split into 3")
						continue
					}

					s := statusItem{
						Service: content[0],
						Version: content[1],
						Health:  content[2],
					}
					dl.Services = append(dl.Services, s)
				}
			}
		}
		log.Debug("finished")

		res, _ := json.Marshal(dl)
		w.Write(res)

	})

	m.HandleFunc("/api/almagest/control/", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-type", "application/json")

		log.Debug("got a control message", "path", r.URL.Path)
		log.Debug("type of message", "control", r.URL.Path[22:])
		if strings.Contain(r.URL.Path[22:], "debug") {
			log.Debug("got a debug type message")
			// publish a control message
			// save it to redis for others to see
		}

		w.Write([]byte("a response"))

	})
	port := "0.0.0.0:39788"

	fmt.Println(gitCommit)
	log.Info("Starting server", "listening", port, "version", strings.TrimSpace(gitCommit))
	http.ListenAndServe(port, m)

}
