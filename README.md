This is a simple TOTP _(Time-based One-time Password)_ CLI tool. You can manage your
accounts with namespaces and protect your data with a PIN _(aka password)_.

```
❯ ./bin/totp-cli
           add-token     Add new token
                         `totp-cli add-token`
                         This command will ask for the namespace, the account and the token
                list     List all available namespaces or accounts under a namespace
                         `totp-cli list`      => list all namespaces
                         `totp-cli list myns` => list all accounts under 'myns' namespace
            generate     Generate a specific OTP
                         `totp-cli generate namespace.account`
                help     This help message :)
```

TODO: Documentation :D
