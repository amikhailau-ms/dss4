package main

import (
	"fmt"

	"github.com/amikhailau/dss4/elliptic_crypto"
)

func main() {
	curve := elliptic_crypto.BuildCurve()
	keysA := elliptic_crypto.Keys{}
	pkAX, pkAY := keysA.GeneratePublicKey(58, curve)
	fmt.Println(pkAX, pkAY)
	keysB := elliptic_crypto.Keys{}
	pkBX, pkBY := keysB.GeneratePublicKey(42, curve)
	fmt.Println(pkBX, pkBY)
	keysA.GeneratePrivateKey(pkBX, pkBY)
	keysB.GeneratePrivateKey(pkAX, pkAY)
	fmt.Println(keysA.ComparePrivateKeys(&keysB))
	message := "Hello, world!"
	r, s := keysA.SignMessage(message)
	fmt.Println(keysB.CheckSignature(message, r, s, pkAX, pkAY))
}
