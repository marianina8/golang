package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the company name (without spaces or special characters): ")
	company, _ := reader.ReadString('\n')

	company = strings.Trim(company, " ")
	company = strings.Trim(company, "\n")

	fmt.Println("Creating account for " + company)

	client, session, err := connectToHost("root", "aspera.dnsname.net:22")
	if err != nil {
		panic(err)
	}

	pass := randStr(10, "alphanum")
	user := "p-" + strings.ToLower(company)

	command := "cd '/Aspera/Incoming files/';"
	command += "mkdir " + company + ";"
	command += "chmod 777 " + company + ";"
	command += "useradd " + user + "; echo '" + pass + "' | passwd " + user + " --stdin;"
	command += "htpasswd -b /opt/aspera/etc/webpasswd " + user + " " + pass + ";"
	command += "asconfigurator -F \"set_user_data;user_name," + user + ";absolute,/Aspera/Incoming files/" + company + "\""

	fmt.Println(command)

	out, err := session.CombinedOutput(command)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	client.Close()

	fmt.Print("\n")
	fmt.Print("\n")
	fmt.Print("\n")
	fmt.Println("Please copy/paste the following information to send to client:")
	fmt.Print("\n")
	fmt.Print("\n")
	fmt.Println("Aspera Credentials:")
	fmt.Println("    username: " + user)
	fmt.Println("    password: " + pass)
	fmt.Println()
	fmt.Println("Link: http://aspera.participant.net/aspera/user")
	fmt.Println()
	fmt.Println("Recommended Browser: ")
	fmt.Println("    Safari on Mac")
	fmt.Println("    Internet Explorer on Windows")
	fmt.Print("\n")
	fmt.Print("\n")

}

func connectToHost(user, host string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password("#########")},
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
