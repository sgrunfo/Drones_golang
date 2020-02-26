package archive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//request body
type RequestBody struct {
	IdShipment string `json:"idShipment"`
	Token      string `json:"token"`
}

//request body percentage
type ReqBodyPerc struct {
	IdShipment string `json:"idShipment"`
	Percentage string `json:"percentage"`
}

//msg to ca
type MessageCa struct {
	Id         string `json:"id"`
	IdShipment string `json:"idShipment"`
	Token      string `json:"token"`
}

// User rappresenta un utente dato il suo Token
type User struct {
	Token string
}

// Shipment rappresenta un pacco dal suo IdTrack e l'utente che l'ha richiesto
type Drone struct {
	IDDrone string
}

func xor(x []byte, y []byte) []byte {
	z := make([]byte, len(x))
	for i := range x {
		z[i] = x[i] ^ y[i]
	}
	return z
}

func getShipmentFile(idShipment string, token string) ([]byte, error) {

	droneFile, err := ioutil.ReadFile("./database/drones.json")
	typeToken := ""
	//var key1, key2, keyReturn []byte

	if err != nil {
		log.Fatal(err)
	}

	var drones []Drone

	err = json.Unmarshal(droneFile, &drones)

	if err != nil {
		fmt.Println(err)
	}

	for _, drone := range drones {
		if drone.IDDrone == token {
			typeToken = "drone"
		}
	}

	usersFile, _ := ioutil.ReadFile("./database/users.json")

	if err != nil {
		log.Fatal(err)
	}

	var users []User

	err = json.Unmarshal(usersFile, &users)

	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.Token == token {
			typeToken = "user"
		}
	}

	if typeToken == "" {
		return nil, errors.New("Token errato")
	} else if typeToken == "drone" {
		shipmentForDrone, _ := ioutil.ReadFile("./archive/shipmentsForDrones/trackId_" + idShipment + ".txt")
		return shipmentForDrone, nil
	} else if typeToken == "user" {
		shipmentsForUser, _ := ioutil.ReadFile("./archive/shipmentsForUsers/id_" + token + "/trackId_" + idShipment + ".txt")
		return shipmentsForUser, nil
	}

	return nil, errors.New("shipment inesistente")
}

func getShipment(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	var msg RequestBody
	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
	}

	file, err1 := getShipmentFile(msg.IdShipment, msg.Token)

	if file != nil {
	}

	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err1.Error()))
	}

	var msgCa MessageCa
	msgCa.Id = msg.Token
	msgCa.IdShipment = msg.IdShipment
	msgCa.Token = "archiveToken"
	msgCaValue, _ := json.Marshal(msgCa)

	resp, errPost := http.Post("http://localhost:1234/getKey", "application/json", bytes.NewBuffer(msgCaValue))

	if errPost != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errPost.Error()))
	}

	key, _ := ioutil.ReadAll(resp.Body)

	msgXor := xor(file, key)

	w.Write(msgXor)
}

func putState(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var msg ReqBodyPerc
	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
	}

	root := "./archive/percentage/"
	search := "percentage_" + msg.IdShipment + ".txt"
	var pathFile string

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if filepath.Ext(path) == ".txt" && info.Name() == search {
			pathFile = path
		}

		return nil
	})

	f, err := os.Create("./" + pathFile)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	//ioutil.WriteFile("./archive/percentage/percentage_"+msg.IdShipment+".json", string(msg.Percentage[:]), 0644)

	f.WriteString(string(msg.Percentage))

}

func getState(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var msg RequestBody
	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
	}

	root := "./archive/percentage/"
	search := "percentage_" + msg.IdShipment + ".txt"
	var pathFile string

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if filepath.Ext(path) == ".txt" && info.Name() == search {
			pathFile = path
		}

		return nil
	})

	file, err1 := ioutil.ReadFile("./archive/percentage/id_" + msg.Token + "/percentage_" + msg.IdShipment + ".txt")

	if file != nil {
	}

	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err1.Error()))
	}

	var msgPerc ReqBodyPerc
	msgPerc.IdShipment = msg.IdShipment
	msgPerc.Percentage = string(file[:])
	msgPercValue, _ := json.Marshal(msgPerc.Percentage)

	w.Write(msgPercValue)

}

func Archive() {
	http.HandleFunc("/getShipment", getShipment)
	http.HandleFunc("/putState", putState)
	http.HandleFunc("/getState", getState)
	//sendKey("111111", "lorenzo", 10)
	//key := getKey("claudio", "idTrackClaudio", "archiveToken")

	var port = "5678"
	fmt.Println("listenig on port " + port + " ...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}

	//key := getKey("claudio", "idTrackClaudio", "drone_2")

}
