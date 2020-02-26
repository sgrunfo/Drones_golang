package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	keys "../keys"
)

// Drone rappresenta un drone
type Drone struct {
	IDDrone string
}

// Shipment rappresenta una singola spedizione
type Shipment struct {
	IdShipment string
}

//msg to ca
type MessageShipment struct {
	IdShipment string `json:"idShipment"`
	Token      string `json:"token"`
}

//msg to ca
type MessageCa struct {
	Id         string `json:"id"`
	IdShipment string `json:"idShipment"`
	Token      string `json:"token"`
}

//msg to ca
type MessagePercentage struct {
	IdShipment string `json:"idShipment"`
	Percentage string `json:"percentage"`
}

// User rappresenta un utente dato il suo Token
type User struct {
	Token string
}

type ShipmentComplete struct {
	IDTrack            string
	User               User
	Source             string
	Target             string
	Altitude           string
	Speed              string
	DepartureTime      string
	ArrivalTime        string
	Weight             string
	Price              string
	CompletePercentage string
}

var privateKey keys.PrivateKey

func xor(x []byte, y []byte) []byte {
	z := make([]byte, len(x))
	for i := range x {
		z[i] = x[i] ^ y[i]
	}
	return z
}

func decriptShipment(id string, idShipment string, token string) []byte {

	//richiesta all'archivio
	var msgShipment MessageShipment
	msgShipment.IdShipment = idShipment
	msgShipment.Token = token
	msgShipmentValue, _ := json.Marshal(msgShipment)
	shipmentResp, err := http.Post("http://localhost:5678/getShipment", "application/json", bytes.NewBuffer(msgShipmentValue))

	if err != nil {
		fmt.Println(err)
	}

	shipment, _ := ioutil.ReadAll(shipmentResp.Body)

	//richiesta alla ca
	var msgCa MessageCa
	msgCa.Id = id
	msgCa.IdShipment = idShipment
	msgCa.Token = token
	msgCaValue, _ := json.Marshal(msgCa)
	keyResp, err1 := http.Post("http://localhost:1234/getKey", "application/json", bytes.NewBuffer(msgCaValue))

	if err1 != nil {
		fmt.Println(err1)
	}

	cryptedKey, _ := ioutil.ReadAll(keyResp.Body)

	z := new(big.Int)
	z.SetBytes(cryptedKey)
	bigKey := keys.Decrypt(z, privateKey)

	msgXor := xor(shipment, bigKey.Bytes())
	return msgXor

}

// definisce cosa può fare il drone
func main() {

	fmt.Println("inserisci l'id del drone: ")
	reader := bufio.NewReader(os.Stdin)
	var id string
	fmt.Print("-> ")
	id, _ = reader.ReadString('\n')
	// convert CRLF to LF
	id = strings.Replace(id, "\n", "", -1)
	id = strings.Replace(id, "\r", "", -1)
	privateKey = keys.KeyPair(id)

	reader = bufio.NewReader(os.Stdin)
	fmt.Println("inserisci l'id della spedizione: ")
	var idShipment string
	fmt.Print("-> ")
	idShipment, _ = reader.ReadString('\n')
	// convert CRLF to LF
	idShipment = strings.Replace(idShipment, "\n", "", -1)
	idShipment = strings.Replace(idShipment, "\r", "", -1)

	shipment := decriptShipment(id, idShipment, id)
	var shipmentJson ShipmentComplete
	err2 := json.Unmarshal(shipment, &shipmentJson)
	if err2 != nil {
		fmt.Println(err2)
	}
	var messagePercentage MessagePercentage
	messagePercentage.IdShipment = idShipment
	var complete = 0

	for complete < 100 {
		time.Sleep(15 * time.Second)

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		percentage := complete + r1.Intn(40)
		if percentage > 100 {
			complete = 100
		} else {
			complete = percentage
		}
		fmt.Println("completamento: " + strconv.Itoa(complete) + "%")
		messagePercentage.Percentage = strconv.Itoa(complete)
		msgPercentageValue, _ := json.Marshal(messagePercentage)
		_, err1 := http.Post("http://localhost:5678/putState", "application/json", bytes.NewBuffer(msgPercentageValue))
		if err1 != nil {
			fmt.Println(err1)
		}
	}
	shipment = nil
	_ = json.Unmarshal([]byte(""), &shipmentJson)

}
