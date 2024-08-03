package main

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

func main() {
	password := "xReJjM3fhY6cFumZiKUYykW05htave5wfE5tU0cryLLEseeaUbis7UdRYyhiraOa"
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))

	fmt.Println("Encoded Hash:", string(encoded))
	if err != nil {
		panic(err)
	}

	// Verify the password
	match, err := argon2.VerifyEncoded([]byte(password), encoded)
	if err != nil {
		fmt.Println("Error verifying password:", err)
		return
	}
	if match {
		fmt.Println("Password is correct!")
	} else {
		fmt.Println("Password is incorrect!")
	}
	match, err = argon2.VerifyEncoded([]byte(password), []byte("$argon2id$v=19$m=65540,t=3,p=4$SmQveXFYdUtQY1NFUkQ5elJYWXI1MGc4WkljRW1xUVJPVGhrRHlZL3hUUT0$evQ5G6+c+PNw0dOhKYcyvXttOkcPtZwpCF3bnHcOePI"))
	if err != nil {
		fmt.Println("Error verifying password:", err)
		return
	}
	if match {
		fmt.Println("Password is correct!")
	} else {
		fmt.Println("Password is incorrect!")
	}
	match, err = argon2.VerifyEncoded([]byte("wrong"), encoded)
	if err != nil {
		fmt.Println("Error verifying password:", err)
		return
	}
	if match {
		fmt.Println("Password is correct!")
	} else {
		fmt.Println("Password is incorrect!")
	}

	match, err = argon2.VerifyEncoded([]byte("test"), []byte("$argon2id$v=19$m=65536,t=3,p=4$DQZdkfMaVVZCnm/GyxZCPA$dNPCNuwFa7J7HY+gB+O096y0hJfzzpo+iC8hRBWmLv8"))
	if err != nil {
		fmt.Println("Error verifying password:", err)
		return
	}
	if match {
		fmt.Println("Password is correct!")
	} else {
		fmt.Println("Password is incorrect!")
	}
}
