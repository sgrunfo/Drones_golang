package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"

	keys "../keys"
)

//request body percentage shipment
type RequestPerc struct {
	IdShipment string `json:"idShipment"`
	Token      string `json:"token"`
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

func getPercentage(idShipment string, token string) string {

	var msgPerc RequestPerc
	msgPerc.IdShipment = idShipment
	msgPerc.Token = token
	msgPercValue, _ := json.Marshal(msgPerc)
	percResp, err := http.Post("http://localhost:5678/getState", "application/json", bytes.NewBuffer(msgPercValue))

	if err != nil {
		fmt.Println(err)
	}

	percentage, _ := ioutil.ReadAll(percResp.Body)

	return string(percentage[:])

}

// User definisce cosa può fare l'utente
func main() {

	fmt.Println("inserisci il tuo id: ")
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

	perc := getPercentage(idShipment, id)

	fmt.Println("----------------------------------------")
	fmt.Println("Id spedizione: " + shipmentJson.IDTrack)
	fmt.Println("Arrivo previsto: " + shipmentJson.ArrivalTime)
	fmt.Println("Arrivo: " + shipmentJson.Target)
	fmt.Println("Prezzo: " + shipmentJson.Price)
	fmt.Println("Percentuale di completamento: " + perc)
	fmt.Println("----------------------------------------")

	for {
		reader = bufio.NewReader(os.Stdin)
		fmt.Println("inserisci: 'state' per ricevere lo stato della spedizione, 'exit' per terminare ")
		var perc string
		fmt.Print("-> ")
		perc, _ = reader.ReadString('\n')
		perc = strings.Replace(perc, "\n", "", -1)
		perc = strings.Replace(perc, "\r", "", -1)
		if strings.Compare("state", perc) == 0 {
			perc := getPercentage(idShipment, id)
			fmt.Println("----------------------------------------")
			fmt.Println("Percentuale di completamento: " + perc)
			fmt.Println("----------------------------------------")
		} else if strings.Compare("exit", perc) == 0 {
			fmt.Println("Ciao!")
			return
		}
	}

}
