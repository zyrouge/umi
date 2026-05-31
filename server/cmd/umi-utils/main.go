package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"zyrouge.me/umi/utils"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	var err error = nil
	switch cmd {
	case "help":
		printUsage()
		break
	case "hash-password":
		err = hashPassword()
	case "generate-hex-key":
		err = generateHexKey()
	case "generate-base64-key":
		err = generateBase64Key()
	default:
		err = fmt.Errorf("unknown command: %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("usage: umi-utils <command> [arguments]")
	fmt.Println("commands:")
	fmt.Println("  hash-password <password>  Hash a password using bcrypt")
	fmt.Println("  generate-hex-key          Generate a random 32-byte key and output as hex (jwt_secret)")
	fmt.Println("  generate-base64-key       Generate a random 32-byte key and output as base64 (team_encryption_key and user_encryption_key)")
}

func hashPassword() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("invalid command usage")
	}
	password := os.Args[2]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	fmt.Println(string(hash))
	return nil
}

func generateHexKey() error {
	bytes, err := utils.GenerateRandomBytes(32)
	if err != nil {
		return err
	}
	fmt.Println(utils.BytesToHex(bytes))
	return nil
}

func generateBase64Key() error {
	bytes, err := utils.GenerateRandomBytes(32)
	if err != nil {
		return err
	}
	fmt.Println(utils.BytesToBase64(bytes))
	return nil
}
