package cmd

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "net/url"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "github.com/urfave/cli/v2"
    s "github.com/yitsushi/totp-cli/internal/storage"
    "github.com/yitsushi/totp-cli/internal/terminal"
)

func promptForInput(prompt, defaultValue string) string {
    reader := bufio.NewReader(os.Stdin)
    if defaultValue != "" {
        prompt = fmt.Sprintf("%s (%s)", prompt, defaultValue)
    }
    fmt.Print(prompt + ": ")
    text, err := reader.ReadString('\n')
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
        os.Exit(1)
    }
    text = strings.TrimSpace(text)
    if text == "" {
        text = defaultValue
    }
    return text
}


func installZbarimg() error {
    var cmd *exec.Cmd

    switch os := runtime.GOOS; os {
    case "linux":
        cmd = exec.Command("sudo", "apt-get", "install", "zbar-tools")
    case "darwin":
        cmd = exec.Command("brew", "install", "zbar")
    default:
        return fmt.Errorf("unsupported platform: %s", os)
    }

    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("failed to install zbarimg: %w", err)
    }

    return nil
}

func processQRCode() (string, error) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the path to the QR code file or directory: ")
    path, err := reader.ReadString('\n')
    if err != nil {
        return "", fmt.Errorf("failed to read input: %w", err)
    }

    path = strings.TrimSpace(path)

    info, err := os.Stat(path)
    if err != nil {
        return "", fmt.Errorf("failed to get file or directory info: %w", err)
    }

    var files []string
    if info.IsDir() {
        entries, err := ioutil.ReadDir(path)
        if err != nil {
            return "", fmt.Errorf("failed to read directory: %w", err)
        }

        for _, entry := range entries {
            if !entry.IsDir() {
                files = append(files, filepath.Join(path, entry.Name()))
            }
        }
    } else {
        files = append(files, path)
    }

    tempFile, err := ioutil.TempFile("", "qr-output-*.txt")
    if err != nil {
        return "", fmt.Errorf("failed to create temporary file: %w", err)
    }

    for _, file := range files {
		cmd := exec.Command("zbarimg", "--raw", "-q", file)

		output, err := cmd.Output()
		if err != nil {
			fmt.Println("It seems zbarimg isn't installed. Installing...")
			installErr := installZbarimg()
			if installErr != nil {
				fmt.Println("Automatic installation failed. Please install zbarimg manually to use the --use-qr functionality.")
				return "", fmt.Errorf("failed to install zbarimg: %w", installErr)
			}

			cmd := exec.Command("zbarimg", "--raw", "-q", file)
			output, err = cmd.Output()
			if err != nil {
				fmt.Println("is your file path correct? Something went wrong.")
				return "", fmt.Errorf("failed to run zbarimg: %w", err)
			}
		}

        _, err = tempFile.Write(output)
        if err != nil {
            return "", fmt.Errorf("failed to write to temporary file: %w", err)
        }
    }

    return tempFile.Name(), nil
}

func processOTPAuth(filePath string, silq bool) (string, string, string, error) {
	// useSilq := silq
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", "", "", fmt.Errorf("failed to read file: %w", err)
    }
    lines := strings.Split(string(content), "\n")
    var otpauthLines []string
    for _, line := range lines {
        if strings.HasPrefix(line, "otpauth") {
            otpauthLines = append(otpauthLines, line)
        }
    }
    var namespace, account, token string
    if len(otpauthLines) > 0 {
        otpauthURL, err := url.Parse(otpauthLines[0])
        if err != nil {
            return "", "", "", fmt.Errorf("failed to parse otpauth link: %w", err)
        }

        pathParts := strings.Split(otpauthURL.Path, "/")
        if len(pathParts) >= 3 {
            // namespace = pathParts[1]
            account = pathParts[2]
        }

		namespace = otpauthURL.Query().Get("issuer")
        if namespace == "" && silq {
            namespace = "qrimports"
            fmt.Println("No issuer found in otpauth link. Using default namespace 'qrimports'.")
        }

        token = otpauthURL.Query().Get("secret")
    }
    if !silq {
        namespace = promptForInput("Namespace:", namespace)
    }


    return namespace, account, token, nil
}


// ... existing code ...
func AddTokenCommand() *cli.Command {
    return &cli.Command{
		Name:      "add-token",
		Aliases:   []string{"add"},
		Usage:     "Add new token.",
		ArgsUsage: "[namespace] [account]",

        Flags: []cli.Flag{
			&cli.UintFlag{
				Name:  "length",
				Value: s.DefaultTokenLength,
				Usage: "Length of the generated token.",
			},
			&cli.StringFlag{
				Name:  "prefix",
				Value: "",
				Usage: "Prefix for the token.",
			},

            &cli.BoolFlag{
                Name:  "use-qr",
                Usage: "Use QR code for token generation.",
            },
            &cli.BoolFlag{
                Name:  "use-qrl",
                Usage: "Use QR code link for token generation.",
            },
			&cli.BoolFlag{
				Name:  "silq",
				Usage: "Use the namespace from the otpauth link.",
			},
        },
        Action: func(ctx *cli.Context) error {
			silq := ctx.Bool("silq")

			var (
				namespace *s.Namespace
				account   *s.Account
				err       error
			)

			if ctx.Bool("use-qr") {			

				tempFilePath, err := processQRCode()
				if err != nil {
					return err
				}
				
				nsName, accName, token, err := processOTPAuth(tempFilePath, silq)
				if err != nil {
					return err
				}
				
				err = os.Remove(tempFilePath)
				if err != nil {
					fmt.Println("Failed to delete temporary file:", err)
					return err
				}
				
				storage, err := s.PrepareStorage()
				if err != nil {
					return err
				}
				
				namespace, err := storage.FindNamespace(nsName)
				if err != nil {
					namespace = &s.Namespace{Name: nsName}
					storage.Namespaces = append(storage.Namespaces, namespace)
				}
				
				account, err := namespace.FindAccount(accName)
				if err == nil {
					return CommandError{
						Message: fmt.Sprintf("%s.%s exists", namespace.Name, account.Name),
					}
				}
				
				account = &s.Account{Name: accName, Token: token, Prefix: ctx.String("prefix"), Length: ctx.Uint("length")}
				namespace.Accounts = append(namespace.Accounts, account)
				
				err = storage.Save()
				if err != nil {
					return err
				}
				
				return nil

			}
			
			if ctx.Bool("use-qrl") {

				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter the path to the link file: ")
				filePath, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Failed to read input.")
					return err
				}
				
				filePath = strings.TrimSpace(filePath)
				
				nsName, accName, token, err := processOTPAuth(filePath, silq)
				if err != nil {
					return fmt.Errorf("failed to process otpauth: %w", err)
				}
				
				storage, err := s.PrepareStorage()
				if err != nil {
					return err
				}
				
				namespace, err := storage.FindNamespace(nsName)
				if err != nil {
					namespace = &s.Namespace{Name: nsName}
					storage.Namespaces = append(storage.Namespaces, namespace)
				}
				
				account, err := namespace.FindAccount(accName)
				if err == nil {
					return CommandError{
						Message: fmt.Sprintf("%s.%s exists", namespace.Name, account.Name),
					}
				}
				
				account = &s.Account{Name: accName, Token: token, Prefix: ctx.String("prefix"), Length: ctx.Uint("length")}
				namespace.Accounts = append(namespace.Accounts, account)
				
				err = storage.Save()
				if err != nil {
					return err
				}
				
				return nil
			}

				nsName, accName, token := askForAddTokenDetails(
					ctx.Args().Get(argSetPrefixPositionNamespace),
					ctx.Args().Get(argSetPrefixPositionAccount),
				)
		
				storage, err := s.PrepareStorage()
				if err != nil {
					return err
				}


				namespace, err = storage.FindNamespace(nsName)
					if err != nil {
						namespace = &s.Namespace{Name: nsName}
						storage.Namespaces = append(storage.Namespaces, namespace)
					}


					account, err = namespace.FindAccount(accName)
					if err == nil {
						return CommandError{
							Message: fmt.Sprintf("%s.%s exists", namespace.Name, account.Name),
						}
					}
					
					account = &s.Account{Name: accName, Token: token, Prefix: ctx.String("prefix"), Length: ctx.Uint("length")}
					namespace.Accounts = append(namespace.Accounts, account)

					err = storage.Save()
						if err != nil {
							return err
						}
						return nil


    	},
	}
}


func askForAddTokenDetails(namespace, account string) (string, string, string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	token, _ := term.Read("Token:")

	return namespace, account, token
}
