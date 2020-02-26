package ca

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"time"

	keys "../keys"
)

// User rappresenta un utente dato il suo Token
type User struct {
	Token string
}

// Drone rappresenta un drone
type Drone struct {
	IDDrone string
}

// ShipmentComplete rappresenta tutte le informazioni del pacco
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

type Message struct {
	Id         string `json:"id"`
	IdShipment string `json:"idShipment"`
	Token      string `json:"token"`
}

var privateKey = keys.KeyPair("CA")

func createKey(length int) []byte {
	//fmt.Println("CA: Sto generando una nuova chiave")
	key := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(key)
	return key
}

func createKeySeed(length int, seedSource int) []byte {
	//fmt.Println("CA: Sto generando un nuovo seed")
	key := make([]byte, length)
	r := rand.New(rand.NewSource(int64(seedSource)))
	r.Read(key)
	return key
}

func xor(x []byte, y []byte) []byte {
	z := make([]byte, len(x))
	for i := range x {
		z[i] = x[i] ^ y[i]
	}
	return z
}

func saveID(id string) {
	os.Mkdir("./ca/ids/id_"+id, os.ModePerm)
	//fmt.Println("CA: ho salvato id: " + id)
}

func saveShipment(idShipment string, idUser string) {
	os.Mkdir("./ca/ids/id_"+idUser+"/sm_"+idShipment, os.ModePerm)
	//fmt.Println("CA: ho salvato idShipment: " + idShipment)
	//fmt.Println("CA: ho salvato isUser: " + idUser)
}

func saveKey(idShipment string, idUser string, key1 []byte, key2 []byte) {
	//fmt.Println("idShipment: " + idShipment)
	//fmt.Println("isUser: " + idUser)

	ioutil.WriteFile("./ca/ids/id_"+idUser+"/sm_"+idShipment+"/key1/key1", key1, 0644)
	//fmt.Println("CA: ho salvato la Key1")
	ioutil.WriteFile("./ca/ids/id_"+idUser+"/sm_"+idShipment+"/key2/key2", key2, 0644)
	//fmt.Println("CA: ho salvato la Key2")
}

// func sendKey(idShipment string, idUser string, length int) {
// 	saveID(idUser) //invece di creare la cartella controlla solo se esiste l'utente altrimenti da errore
// 	saveShipment(idShipment, idUser)

// 	key1 := createKey(length)
// 	key2 := createKey(length)
// 	key2 = xor(key1, key2)
// 	os.Mkdir("./ca/ids/id_"+idUser+"/sm_"+idShipment+"/key1", os.ModePerm)
// 	os.Mkdir("./ca/ids/id_"+idUser+"/sm_"+idShipment+"/key2", os.ModePerm)

// 	saveKey(idShipment, idUser, key1, key2)
// }

// Initialize crea le chiavi e salva per la CA
func InitializeKeysForUsers(idShipment string, idUser string, lenght int) {
	os.Mkdir("./ca/keysForUsers", os.ModePerm)
	os.Mkdir("./ca/keysForUsers/id_"+idUser, os.ModePerm)
	os.Mkdir("./ca/keysForUsers/id_"+idUser+"/trackId_"+idShipment, os.ModePerm)

	key1 := createKey(lenght)
	//fmt.Println("CA: ho creato la Key1")
	data := int(binary.BigEndian.Uint64(key1))
	key2 := createKeySeed(lenght, data+int(time.Now().UnixNano()))
	//fmt.Println("CA: ho creato la Key2")
	ioutil.WriteFile("./ca/keysForUsers/id_"+idUser+"/trackId_"+idShipment+"/key1.txt", key1, 0644)
	ioutil.WriteFile("./ca/keysForUsers/id_"+idUser+"/trackId_"+idShipment+"/key2.txt", key2, 0644)

}

// Initialize crea le chiavi e salva per la CA
func InitializeKeysForDrones(idShipment string, lenght int) {
	os.Mkdir("./ca/keysForDrones", os.ModePerm)
	os.Mkdir("./ca/keysForDrones/trackId_"+idShipment, os.ModePerm)

	key1 := createKey(lenght)
	//fmt.Println("CA: ho creato la Key1")
	data := int(binary.BigEndian.Uint64(key1))
	key2 := createKeySeed(lenght, data+int(time.Now().UnixNano()))
	//fmt.Println("CA: ho creato la Key2")
	ioutil.WriteFile("./ca/keysForDrones/trackId_"+idShipment+"/key1.txt", key1, 0644)
	ioutil.WriteFile("./ca/keysForDrones/trackId_"+idShipment+"/key2.txt", key2, 0644)

}

//idReq id relativo a chi serve decriptare la spedizione
//idShipment identificativo spedizione
//token di chi fa la richiesta

func getKey(idReq string, idShipment string, token string) []byte {

	droneFile, err := ioutil.ReadFile("./database/drones.json")
	typeReq := ""
	typeToken := ""
	var key1, key2, keyReturn []byte

	if err != nil {
		log.Fatal(err)
	}

	var drones []Drone

	err = json.Unmarshal(droneFile, &drones)

	if err != nil {
		fmt.Println(err)
	}

	for _, drone := range drones {
		if drone.IDDrone == idReq {
			typeReq = "drone"
		}
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
		if user.Token == idReq {
			typeReq = "user"
		}
		if user.Token == token {
			typeToken = "user"
		}
	}

	if typeReq == "" && typeToken == typeToken && ((keyReturn != nil) || (keyReturn == nil)) {
		fmt.Println(idReq)
		fmt.Println(keyReturn)
		fmt.Println("id non corrisponde a nessuno")
		return nil
	}

	if token == "archiveToken" {
		if typeReq == "user" {
			keyReturn, _ = ioutil.ReadFile("./ca/keysForUsers/id_" + idReq + "/trackId_" + idShipment + "/key2.txt")
		} else if typeReq == "drone" {
			keyReturn, _ = ioutil.ReadFile("./ca/keysForDrones/trackId_" + idShipment + "/key2.txt")
		}
	} else {

		if idReq != token {
			fmt.Println("Non sei abilitato a fare questa richiesta")
			return nil
		}

		if typeReq == "user" {
			key1, _ = ioutil.ReadFile("./ca/keysForUsers/id_" + idReq + "/trackId_" + idShipment + "/key1.txt")
			key2, _ = ioutil.ReadFile("./ca/keysForUsers/id_" + idReq + "/trackId_" + idShipment + "/key2.txt")
			xoredKeys := xor(key1, key2)

			userKey := keys.GetPublicKey(token)
			z := new(big.Int)
			z.SetBytes(xoredKeys)

			bigIntKey := keys.Encrypt(z, userKey)
			keyReturn = bigIntKey.Bytes()
		} else if typeReq == "drone" {
			key1, _ = ioutil.ReadFile("./ca/keysForDrones/trackId_" + idShipment + "/key1.txt")
			key2, _ = ioutil.ReadFile("./ca/keysForDrones/trackId_" + idShipment + "/key2.txt")
			xoredKeys := xor(key1, key2)

			droneKey := keys.GetPublicKey(token)
			z := new(big.Int)
			z.SetBytes(xoredKeys)

			bigIntKey := keys.Encrypt(z, droneKey)
			keyReturn = bigIntKey.Bytes()
		}
	}

	return keyReturn
}

func getKeyHttp(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	var msg Message
	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
	}
	key := getKey(msg.Id, msg.IdShipment, msg.Token)

	w.Write(key)
}

func Ca() {
	http.HandleFunc("/getKey", getKeyHttp)
	//sendKey("111111", "lorenzo", 10)
	//key := getKey("claudio", "idTrackClaudio", "archiveToken")
	var port = "1234"
	fmt.Println("listenig on port " + port + " ...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}

	//key := getKey("claudio", "idTrackClaudio", "drone_2")

}
