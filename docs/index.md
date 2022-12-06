[![Go Report Card](https://goreportcard.com/badge/github.com/yitsushi/totp-cli)](https://goreportcard.com/report/github.com/yitsushi/totp-cli)
[![Actions Status](https://github.com/yitsushi/totp-cli/actions/workflows/quality-check.yaml/badge.svg)](https://github.com/yitsushi/totp-cli/actions/workflows/quality-check.yaml)
[![Coverage Status](https://coveralls.io/repos/github/yitsushi/totp-cli/badge.svg?branch=main)](https://coveralls.io/github/yitsushi/totp-cli?branch=main)

This is a simple TOTP _(Time-based One-time Password)_ CLI tool.
TOTP is the most common mechanism for 2FA _(Two-Factor-Authentication)_.
You can manage and organize your accounts with namespaces
and protect your data with a password.

### Install

Download the latest version of the application
from the [releases page](https://github.com/yitsushi/totp-cli/releases/latest).

Users on macOS can also install the package using [MacPorts](https://ports.macports.org/port/totp-cli/summary):

```
sudo port selfupdate
sudo port install totp-cli
```

or [Homebrew](https://brew.sh/):

```
brew install totp-cli
```

### Update

```
$ totp-cli update
```

### Help output

```
$ totp-cli help

version                           Print current version of this application
add-token [namespace] [account]   Add new token
delete <namespace> [account]      Delete an account or a whole namespace
instant                           Generate an OTP from TOTP_TOKEN or stdin without the Storage backend
list [namespace]                  List all available namespaces or accounts under a namespace
update                            Check and update totp-cli itself
change-password                   Change password
dump [namespace]                  Dump all available namespaces or accounts under a namespace
generate <namespace> <account>    Generate a specific OTP
help [command]                    Display this help or a command specific help
```

### Usage

When you run the application for the first time, it will ask
for your password. **DO NOT FORGET IT!** There is no way to
recover your password if you forget it.

Your first command _(after `help`)_ would be `add-token`. You get get
your token read a TOTP QR Code.

```
$ totp-cli add-token
Namespace: personal
Account: digitalocean
Token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Password: ***
```

You can specify the namespace and the account name as a parameter:

```
$ totp-cli add-token personal randomaccount
Token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Password: ***
```

If you want to delete `randomaccount` _(because it was a test for example)_,
you can use `delete`:

```
$ totp-cli delete personal.randomaccount
Password: ***
You want to delete 'personal.randomaccount' account.
Are you sure? yes
```

After few accounts, it's a bit hard to remember what did you added,
so you can list namespaces:

```
$ totp-cli list
Password: ***
company1 (Number of accounts: 3)
company2 (Number of accounts: 5)
personal (Number of accounts: 8)
```

or you can list your accounts under a specific namespace:

```
$ totp-cli list personal
Password: ***
personal.evernote
personal.google
personal.github
personal.ifttt
personal.digitalocean
personal.dropbox
personal.facebook
```

If you want to change your password,
you can do it with the `change-password` command.

### Changing the location of the credentials file

Simply put this into your `.zshrc` (or `.{YourShell}rc` or `.profile`):

```
export TOTP_CLI_CREDENTIAL_FILE="/mnt/mydrive/totp-credentials"
```

Or call the client with `TOTP_CLI_CREDENTIAL_FILE`:

```
$ TOTP_CLI_CREDENTIAL_FILE=/mnt/mydrive/totp-credentials totp-cli list
```

Note: It's a filename not just a directory.

Note: It does not traverse through the given path,
      parent directory has to be there already.

### Zsh Completion

A function to provide tab-completion for zsh is in the file `_totp-cli`.
When installing or packaging totp-cli this should preferably be
installed in `$prefix/share/zsh/site-functions`. Otherwise, it can be
installed by copying to a directory where zsh searches for completion
functions (the `$fpath` array). If you, for example, put all completion
functions into the folder `~/.zsh/completions` you must add the
following to your zsh main config file (`.zshrc`):

```
fpath=( ~/.zsh/completions $fpath )
autoload -U compinit
compinit
```
