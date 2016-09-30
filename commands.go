package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type CommandFunction func()

var storage *Storage

var commandHandlers map[string]CommandFunction = map[string]CommandFunction{
	"generate":   Command_Generate,
	"help":       Command_Help,
	"add-token":  Command_AddToken,
	"list":       Command_List,
	"delete":     Command_Delete,
	"change-pin": Command_ChangePIN,
	"update":     Command_Update,
	"version":    Command_Version,
}

var commandDescriptions map[string]string = map[string]string{
	"generate":   "Generate a specific OTP%NLWI%`totp-cli generate namespace.account`",
	"help":       "This help message :)",
	"add-token":  "Add new token%NLWI%`totp-cli add-token`%NLWI%This command will ask for the namespace, the account and the token",
	"list":       "List all available namespaces or accounts under a namespace%NLWI%`totp-cli list`      => list all namespaces%NLWI%`totp-cli list myns` => list all accounts under 'myns' namespace",
	"delete":     "Delete an account or a whole namespace%NLWI%`totp-cli delete nsname`%NLWI%`totp-cli delete nsname.accname`",
	"change-pin": "Change PIN code",
	"update":     "Update totp-cli itself",
	"version":    "Current version number of this application",
}

func prepareStorage() {
	initStorage()

	if storage != nil {
		return
	}

	pin := AskPIN(32, "")

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

	namespace.Accounts = append(namespace.Accounts, account)

	storage.Save()
}

func Command_Delete() {
	term := flag.Arg(1)
	if len(term) < 1 {
		Command_Help()
		return
	}

	path := strings.Split(term, ".")

	nsName := path[0]
	accName := ""

	if len(path) > 1 {
		accName = path[1]
	}

	prepareStorage()

	namespace, err := storage.FindNamespace(nsName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if accName != "" {
		account, err := namespace.FindAccount(accName)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("You want to delete '%s.%s' account.\n", namespace.Name, account.Name)

		if Confirm("Are you sure?") {
			namespace.DeleteAccount(account)
			storage.Save()
		}

	} else {
		fmt.Printf("You want to delete '%s' namespace with %d accounts.\n", namespace.Name, len(namespace.Accounts))
		for _, account := range namespace.Accounts {
			fmt.Printf(" - %s.%s\n", namespace.Name, account.Name)
		}

		if Confirm("Are you sure?") {
			storage.DeleteNamespace(namespace)
			storage.Save()
		}
	}
}

func Command_ChangePIN() {
	prepareStorage()
	newPIN := AskPIN(32, "New PIN")
	newPINConfirm := AskPIN(32, "Again")

	if !CheckPINConfirm(newPIN, newPINConfirm) {
		fmt.Println("New PIN and the confirm mismatch!")
		return
	}

	storage.PIN = newPIN
	storage.Save()
}

func Command_Update() {
	updater := &Updater{}
	if updater.CheckAndDownloadVersion() {
		fmt.Printf("Now you have a fresh new %s \\o/\n", AppName)
	} else {
		fmt.Printf("Your %s is up-to-date. \\o/\n", AppName)
	}
}

func Command_Version() {
	fmt.Printf("%s %s (%s/%s)\n", AppName, AppVersion, runtime.GOOS, runtime.GOARCH)
}

func Command_NotImplementedYet() {
	fmt.Println(" -- Not Implemented Yet --")
}

func askForAddTokenDetails() (namespace, account, token string) {
	namespace = flag.Arg(1)
	account = flag.Arg(2)
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

func initStorage() {
	currentUser, err := user.Current()
	check(err)
	homePath := currentUser.HomeDir
	documentDirectory := filepath.Join(homePath, ".config/totp-cli")

	if _, err := os.Stat(documentDirectory); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(documentDirectory, 0700)
			check(err)
		} else {
			check(err)
		}
	}

	if _, err := os.Stat(filepath.Join(documentDirectory, "credentials")); err == nil {
		return
	}

	pin := AskPIN(32, "You PIN/Password (do not forget it)")
	storage = &Storage{
		File: filepath.Join(documentDirectory, "credentials"),
		PIN:  pin,
	}

	storage.Save()
}
