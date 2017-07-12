package main

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"log"

	"crypto/rsa"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	begin      = []byte("-----BEGIN PUBLIC KEY-----")
	end        = []byte("-----END PUBLIC KEY-----")
	newLine    = []byte("\n")
	lineLength = 64
)

/*
type VerifyKey struct {
	VerifyKey *rsa.PublicKey `json:"verify_key"`
}
*/

func GetVerificationKey() (*rsa.PublicKey, error) {
	appEnv, err := cfenv.Current()
	var vkey string
	if err != nil {
		log.Println(err)
		log.Panic("Error, cf env not available", appEnv)
		panic("Error, cf env not available")
	}

	xsuaaService, err := appEnv.Services.WithName("thingconf-uaa")
	if err != nil {
		log.Println("Error, xsuaa binding not found", err)
		panic("no xsuaa binding")
	}

	var ok bool
	if vkey, ok = xsuaaService.CredentialString("verificationkey"); !ok {
		log.Println("no verification key found")
	}

	var block *pem.Block
	block, _ = pem.Decode(formatPEM([]byte(vkey)))
	if block == nil {
		return nil, fmt.Errorf("pem.Decode failed: %v", err)
	}

	key := formatPEM([]byte(vkey))

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)

	if err != nil {
		return nil, fmt.Errorf("jwt.ParseRSAPublicKeyFromPEM failed: %v", err)
	}

	return publicKey, nil
}

func formatPEM(b []byte) []byte {

	s := bytes.TrimPrefix(b, begin)
	s = bytes.TrimSuffix(s, end)
	pemString := [][]byte{newLine, begin, newLine}
	itr := len(s) / lineLength

	for loop := 0; loop <= itr; loop++ {
		var pos int
		pos = loop * (lineLength)
		if len(s[pos:]) < lineLength {
			tmp := s[pos : pos+len(s[pos:])]
			if pos != pos+len(s[pos:]) {
				pemString = append(pemString, tmp, []byte(newLine))
			}
			break
		}
		tmp := s[pos : pos+lineLength]
		pemString = append(pemString, tmp, []byte(newLine))
	}
	pemString = append(pemString, []byte(end))
	pemByte := bytes.Join(pemString, []byte(""))
	return pemByte
}
