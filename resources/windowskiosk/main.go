package main

import log "github.com/sirupsen/logrus"

func main() {
	c := loadConfig()

	var file string

	credStr := "|" + c.UserName + "|" + c.Password + "\r\n"

	token, err := requestToken()
	if err != nil {
		log.Errorf("couldn't get token %s", err)
	}
	members := getMembers(token)

	for _, m := range members {
		if contains(m.Resources, resource{Name: c.ResourceName}) {
			println(m.Name)
			file += decodeID(m.RFID) + credStr
		}
	}

	writeFile([]byte(file))
}

func contains(s []resource, r resource) bool {
	for _, v := range s {
		if v.Name == r.Name {
			return true
		}
	}

	return false
}
