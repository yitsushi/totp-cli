package main

import (
	"flag"
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type CommandFunction func()

var storage *Storage

var commandHandlers map[string]CommandFunction = map[string]CommandFunction{
	"generate":  Command_Generate,
	"help":      Command_Help,
	"add-token": Command_AddToken,
	"list":      Command_List,
}

var commandDescriptions map[string]string = map[string]string{
	"generate":  "Generate a specific OTP%NLWI%`totp-cli generate namespace.account`",
	"help":      "This help message :)",
	"add-token": "Add new token%NLWI%`totp-cli add-token`%NLWI%This command will ask for the namespace, the account and the token",
	"list":      "List all available namespaces or accounts under a namespace%NLWI%`totp-cli list`      => list all namespaces%NLWI%`totp-cli list myns` => list all accounts under 'myns' namespace",
}

func prepareStorage() {
	pin := AskPIN(32)

	currentUser, err := user.Current()
	check(err)
	homePath := currentUser.HomeDir

	storage = &Storage{
		File: filepath.Join(homePath, ".config/totp-cli/credentials"),
		PIN:  pin,
	}

	storage.Decrypt()
}

func Command_Generate() {
	term := flag.Arg(1)
	if len(term) < 1 {
		Command_Help()
		return
	}

	path := strings.Split(term, ".")

	if len(path) < 2 {
		Command_Help()
		return
	}

	prepareStorage()

	namespace, err := storage.FindNamespace(path[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	account, err := namespace.FindAccount(path[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(GenerateOTPCode(account.Token, time.Now()))
}

func Command_Help() {
	separator := fmt.Sprintf("%5s", "")
	newLineWithIndent := fmt.Sprintf("\n%20s%s", "", separator)
	for command, message := range commandDescriptions {
		message = strings.Replace(message, "%NLWI%", newLineWithIndent, -1)
		fmt.Printf("%20s%s%s\n", command, separator, message)
	}
}

func Command_List() {
	prepareStorage()
	ns := flag.Arg(1)
	if len(ns) < 1 {
		for _, namespace := range storage.Namespaces {
			fmt.Printf("%s (Number of accounts: %d)\n", namespace.Name, len(namespace.Accounts))
		}

		return
	}

	namespace, err := storage.FindNamespace(ns)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, account := range namespace.Accounts {
		fmt.Printf("%s.%s\n", namespace.Name, account.Name)
	}
}

func Command_AddToken() {
	var namespace *Namespace
	var account *Account
	var err error

	nsName, accName, token := askForAddTokenDetails()

	prepareStorage()

	namespace, err = storage.FindNamespace(nsName)
	if err != nil {
		namespace = &Namespace{Name: nsName}
		storage.Namespaces = append(storage.Namespaces, namespace)
	}

	account, err = namespace.FindAccount(accName)
	if err == nil {
		fmt.Println("%s.%s exists!", namespace.Name, account.Name)
	}

	account = &Account{Name: accName, Token: token}

	namespace.Accounts = []*Account{account}

	storage.Save()
}

func Command_NotImplementedYet() {
	fmt.Println(" -- Not Implemented Yet --")
}

func askForAddTokenDetails() (namespace, account, token string) {
	for len(namespace) < 1 {
		namespace = Ask("Namespace")
	}
	for len(account) < 1 {
		account = Ask("Account")
	}
	for len(token) < 1 {
		token = Ask("Token")
	}

	return
}
