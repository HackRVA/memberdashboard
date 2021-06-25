package api

import (
	"io/ioutil"
	"memberserver/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// um3Camera proxies the camera stream from the ultimaker
func (a API) um3Camera(w http.ResponseWriter, req *http.Request) {
	conf, _ := config.Load()

	println("attempting to get ultimaker stream")
	if conf.UM3StreamURL == "" {
		log.Errorf("UM3StreamURL not set")
		http.Error(w, "the server is not setup to retrieve ultimaker camera", http.StatusInternalServerError)
		return
	}

	println("attempting to get ultimaker stream")

	client := &http.Client{}
	req, err := http.NewRequest("GET", conf.UM3StreamURL, nil)

	if err != nil {
		log.Errorf("error creating request to get um3 stream: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("error getting response from um3 stream: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("error parsing response form stream: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
