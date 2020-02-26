package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	ca "../ca"
)

// User rappresenta un utente dato il suo Token
type User struct {
	Token string
}

// Shipment rappresenta un pacco dal suo IdTrack e l'utente che l'ha richiesto
type Shipment struct {
	IDTrack string
	User    User
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

// Drone rappresenta un drone
type Drone struct {
	IDDrone string
}

func initializeTokens() {

	os.Mkdir("./database", os.ModePerm)

	type Users []User

	var users = Users{
		User{
			"claudio",
		},
		User{
			"lorenzo",
		},
		User{
			"daniele",
		},
	}

	usersJSON, _ := json.Marshal(users)

	ioutil.WriteFile("./database/users.json", usersJSON, 0644)
}

func initializeShipments() {

	os.Mkdir("./database", os.ModePerm)

	type Shipments []Shipment

	var shipments = Shipments{
		Shipment{
			"idTrackClaudio",
			User{
				"claudio",
			},
		},
		Shipment{
			"idTrackLorenzo",
			User{
				"lorenzo",
			},
		},
		Shipment{
			"idTrackLorenzo1",
			User{
				"lorenzo",
			},
		},
		Shipment{
			"idTrackDaniele",
			User{
				"daniele",
			},
		},
	}

	shipmentsJSON, _ := json.Marshal(shipments)

	ioutil.WriteFile("./database/shipments.json", shipmentsJSON, 0644)

}

func fillShipments() []ShipmentComplete {

	shipmentFile, err := ioutil.ReadFile("./database/shipments.json")

	if err != nil {
		log.Fatal(err)
	}
	var shipments []ShipmentComplete

	err = json.Unmarshal(shipmentFile, &shipments)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(shipments)
	//fmt.Println(shipments[0].IdTrack)
	for i := range shipments {
		shipments[i].Altitude = "200"
		shipments[i].ArrivalTime = "16:00"
		shipments[i].DepartureTime = "15:35"
		shipments[i].Price = "999.99 euro"
		shipments[i].Source = "41.852784, 12.604829"
		shipments[i].Speed = "40km/h"
		shipments[i].Target = "41.842855, 12.643976"
		shipments[i].Weight = "1.5kg"
		shipments[i].CompletePercentage = "0"
	}
	fmt.Println("Main: ho riempito i pacchi :")
	fmt.Println(shipments)

	return shipments
}

func xor(x []byte, y []byte) []byte {
	z := make([]byte, len(x))
	for i := range x {
		z[i] = x[i] ^ y[i]
	}
	return z
}

func initializeDrones(n int) {
	var drone Drone
	os.Mkdir("./database", os.ModePerm)
	var drones []Drone

	for i := 1; i <= n; i++ {

		drone.IDDrone = "drone_" + strconv.Itoa(i)

		drones = append(drones, drone)
	}

	dronesJSON, _ := json.Marshal(drones)

	ioutil.WriteFile("./database/drones.json", dronesJSON, 0644)
}

//Init is main func
func main() {
	initializeTokens()
	initializeShipments()
	shipments := fillShipments()
	var key1User []byte
	var key1Drone []byte
	os.Mkdir("./archive/percentage", os.ModePerm)
	for i, shipment := range shipments {

		fmt.Println("Main: creo e salvo le spedizioni per gli utenti")
		shipmentForUser := []byte(`{"IDTrack": "` + shipment.IDTrack + `", "Target":"` + shipment.Target + `", "ArrivalTime": "` + shipment.ArrivalTime + `", "Price": "` + shipment.Price + `", "CompletePercentage": "` + shipment.CompletePercentage + `"}`)
		ca.InitializeKeysForUsers(shipments[i].IDTrack, shipments[i].User.Token, len(shipmentForUser))
		key1User, _ = ioutil.ReadFile("./ca/keysForUsers/id_" + shipments[i].User.Token + "/trackId_" + shipments[i].IDTrack + "/key1.txt")
		shipmentXorKey1User := xor(shipmentForUser, key1User)
		os.Mkdir("./archive/shipmentsForUsers/id_"+shipments[i].User.Token, os.ModePerm)
		ioutil.WriteFile("./archive/shipmentsForUsers/id_"+shipments[i].User.Token+"/trackId_"+shipments[i].IDTrack+".txt", shipmentXorKey1User, 0644)

		fmt.Println("Main: creo e salvo le spedizioni per i droni")
		shipmentForDrone := []byte(`{"IDTrack": "` + shipment.IDTrack + `", "Source":"` + shipment.Source + `", "Target":"` + shipment.Target + `", "Altitude":"` + shipment.Altitude + `", "Speed":"` + shipment.Speed + `", "DepartureTime":"` + shipment.DepartureTime + `", "ArrivalTime": "` + shipment.ArrivalTime + `", "Weight": "` + shipment.Weight + `", "CompletePercentage": "` + shipment.CompletePercentage + `"}`)
		ca.InitializeKeysForDrones(shipments[i].IDTrack, len(shipmentForDrone))
		key1Drone, _ = ioutil.ReadFile("./ca/keysForDrones/trackId_" + shipments[i].IDTrack + "/key1.txt")
		shipmentXorKey1Drone := xor(shipmentForDrone, key1Drone)
		ioutil.WriteFile("./archive/shipmentsForDrones/trackId_"+shipments[i].IDTrack+".txt", shipmentXorKey1Drone, 0644)

		fmt.Println("Main: creo e salvo chiavi del pacco: ", shipment)
		//shipmentBytes, _ := json.Marshal(shipment)

		fmt.Println("Main: creo cartella percentuali con file spedizioni: ")
		os.Mkdir("./archive/percentage/id_"+shipments[i].User.Token, os.ModePerm)
		ioutil.WriteFile("./archive/percentage/id_"+shipments[i].User.Token+"/percentage_"+shipments[i].IDTrack+".txt", []byte(""), 0644)

	}

	initializeDrones(len(shipments) - 1)

	/*

		msgLorenzoXorKey1, _ := ioutil.ReadFile("./archive/shipmentsForUsers/id_lorenzo/trackId_idTrackLorenzo1.txt")
		key1User, _ = ioutil.ReadFile("./ca/keysForUsers/id_lorenzo/trackId_idTrackLorenzo1/key1.txt")
		msg := xor(msgLorenzoXorKey1, key1User)
		var shipmentDecoded ShipmentComplete
		err := json.Unmarshal(msg, &shipmentDecoded)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(shipmentDecoded)

	*/

}
