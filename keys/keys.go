package keys

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

//public Key
type PublicKey struct {
	E *big.Int `json:"e"`
	N *big.Int `json:"n"`
}

//Private Key
type PrivateKey struct {
	D *big.Int `json:"d"`
	N *big.Int `json:"n"`
}

func coprimeMod(l, n *big.Int) *big.Int {
	var iGCDn, iGCDl = big.NewInt(1), big.NewInt(1)
	var uno = big.NewInt(1)
	var due = big.NewInt(2)

	for i := iGCDn.Sub(l, uno); i.Cmp(due) >= 0; i.Sub(i, uno) {

		iGCDl.GCD(nil, nil, i, l)

		if iGCDl.Cmp(uno) == 0 {
			return i
		}

	}

	return uno
}

func invMod(l, e *big.Int) *big.Int {
	var i, uno = big.NewInt(1), big.NewInt(1)

	ok := i.ModInverse(e, l)
	if ok != nil {
		//senza i.Add viene i=e
		//i.Add(i,l)
		return i
	}

	return uno
}

func Encrypt(msg *big.Int, publicKey PublicKey) *big.Int {

	msg.Exp(msg, publicKey.E, publicKey.N)
	return msg
}

func Decrypt(encryptedMsg *big.Int, privateKey PrivateKey) *big.Int {

	encryptedMsg.Exp(encryptedMsg, privateKey.D, privateKey.N)
	return encryptedMsg
}

func genPrime(bits int) *big.Int {

	var bigPrime *big.Int
	var a *big.Int
	var err error
	found := false

	for !found {
		bigPrime, err = rand.Prime(rand.Reader, bits)
		a = big.NewInt(1)

		if err != nil {
			fmt.Println(err)
		}

		foundPrime := bigPrime.ProbablyPrime(100)
		mod4bigPrime := *a.Mod(bigPrime, big.NewInt(4))
		mod4cong3 := false

		if mod4bigPrime.Cmp(big.NewInt(3)) == 0 {
			mod4cong3 = true
		}

		if foundPrime && mod4cong3 {
			found = true
		}
	}

	return bigPrime
}

func genKeys() (PrivateKey, PublicKey) {
	var uno = big.NewInt(1)
	var p = genPrime(1024)
	var q = genPrime(1024)

	var n, l, p1, q1 = big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1)

	n.Mul(p, q)
	p1.Sub(p, uno)
	q1.Sub(q, uno)
	l.Mul(p1, q1)

	var publicKey PublicKey
	publicKey.E = coprimeMod(l, n)
	publicKey.N = n

	var privateKey PrivateKey
	privateKey.D = invMod(l, publicKey.E)
	privateKey.N = n

	return privateKey, publicKey
}

//KeyPair initializes private and public key
//need your id
func KeyPair(id string) PrivateKey {
	os.Mkdir("./ca/publics/", os.ModePerm)
	prKey, puKey := genKeys()

	bytes, _ := json.Marshal(puKey)
	ioutil.WriteFile("./ca/publics/id_"+id+".key", bytes, 0644)

	return prKey
}

//GetPublicKey needs id in input
func GetPublicKey(id string) PublicKey {
	pkeyFile, _ := ioutil.ReadFile("./ca/publics/id_" + id + ".key")
	var pkey PublicKey
	err := json.Unmarshal(pkeyFile, &pkey)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	return pkey
}
