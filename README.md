[![Go Report Card](https://goreportcard.com/badge/github.com/yitsushi/totp-cli)](https://goreportcard.com/report/github.com/yitsushi/totp-cli)
[![Actions Status](https://github.com/yitsushi/totp-cli/actions/workflows/quality-check.yaml/badge.svg)](https://github.com/yitsushi/totp-cli/actions/workflows/quality-check.yaml)
[![Coverage Status](https://coveralls.io/repos/github/yitsushi/totp-cli/badge.svg?branch=main)](https://coveralls.io/github/yitsushi/totp-cli?branch=main)

This is a simple TOTP _(Time-based One-time Password)_ CLI tool.
TOTP is the most common mechanism for 2FA _(Two-Factor-Authentication)_.
You can manage and organize your accounts with namespaces
and protect your data with a password.

### Install

Download the latest version of the application
from the [releases page](https://github.com/yitsushi/totp-cli/releases/latest) or using the `go` tool:

```shell
go install github.com/yitsushi/totp-cli@latest
```

Users on macOS can also install the package using [MacPorts](https://ports.macports.org/port/totp-cli/summary):

```shell
sudo port selfupdate
sudo port install totp-cli
```

or [Homebrew](https://brew.sh/):

```
brew install totp-cli
```

### Update

```shell
totp-cli update
```

#### Upgrading from totp-cli v1.2.7 or below

Starting with totp-cli v1.2.8 a [more secure storage
format](https://github.com/FiloSottile/age) is used. The storage will be
upgraded the first time it is written to by totp-cli. You can force this to
occur by running `totp-cli change-password`.

### Help output

```shell
totp-cli help
```

```
version                                     Print current version of this application
delete <namespace> [account]                Delete an account or a whole namespace
dump                                        Dump all available namespaces or accounts under a namespace
instant                                     Generate an OTP from TOTP_TOKEN or stdin without the Storage backend
update                                      Check and update totp-cli itself
list [namespace]                            List all available namespaces or accounts under a namespace
set-prefix [namespace] [account] [prefix]   Set prefix for a token
add-token [namespace] [account]             Add new token
change-password                             Change password
generate <namespace> <account>              Generate a specific OTP
import <input-file>                         Import tokens from a yaml file.
help [command]                              Display this help or a command specific help
```

### Usage

When you run the application for the first time, it will ask
for your password. **DO NOT FORGET IT!** There is no way to
recover your password if you forget it.

Your first command _(after `help`)_ would be `add-token`. You get get
your token read a TOTP QR Code.

```shell
totp-cli add-token
```

```
Namespace: personal
Account: digitalocean
Token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Password: ***
```

You can specify the namespace and the account name as a parameter:

```shell
totp-cli add-token personal randomaccount
```

```
Token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Password: ***
```

If you want to delete `randomaccount` _(because it was a test for example)_,
you can use `delete`:

```shell
totp-cli delete personal.randomaccount
```

```
Password: ***
You want to delete 'personal.randomaccount' account.
Are you sure? yes
```

After few accounts, it's a bit hard to remember what did you added,
so you can list namespaces:

```shell
totp-cli list
```

```
Password: ***
company1 (Number of accounts: 3)
company2 (Number of accounts: 5)
personal (Number of accounts: 8)
```

or you can list your accounts under a specific namespace:

```shell
totp-cli list personal
```

```
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

A prefix can be set with `set-prefix`:

```
totp-cli set-prefix ns account
Prefix: myprefix

# Or with positional argument
totp-cli set-prefix ns account myprefix

# Clear prefix
totp-cli set-prefix ns account --clear
```

To generate an OTP, you simply use the `generate` command like:

```shell
totp-cli generate namespace account
Password: ***
889840
```

You can also use the `--follow` flag on the `generate` command if you want to 
have the OTP token automatically refreshed once expired:

```shell
totp-cli generate --follow namespace account
Password: ***
889840
343555
463346
```


### Changing the location of the credentials file

Simply put this into your `.zshrc` (or `.{YourShell}rc` or `.profile`):

```shell
export TOTP_CLI_CREDENTIAL_FILE="/mnt/mydrive/totp-credentials"
```

Or call the client with `TOTP_CLI_CREDENTIAL_FILE`:

```shell
$ TOTP_CLI_CREDENTIAL_FILE=/mnt/mydrive/totp-credentials totp-cli list
```

The default location is `${HOME}/.config/totp-cli/credentials`.

**Note:** It's a filename not just a directory.

**Note:** It does not traverse through the given path,
      parent directory has to be there already.

### Import

You can import tokens from a YAML file. The syntax is the same as the output of
the `dump` command.

```yaml
- name: ns1
  accounts:
    - name: acc1
      token: updatedtoken
    - name: acc2
      token: mytoken
    - name: acc3
      token: tokenish
- name: ns2
  accounts:
    - name: acc1
      token: token
      prefix: myprefix
```

If a token already exists, it will ask you if you want to overwrite it or not.

```shell
totp-cli import list.yaml
```

### Zsh Completion

A function to provide tab-completion for zsh is in the file `_totp-cli`.
When installing or packaging totp-cli this should preferably be
installed in `$prefix/share/zsh/site-functions`. Otherwise, it can be
installed by copying to a directory where zsh searches for completion
functions (the `$fpath` array). If you, for example, put all completion
functions into the folder `~/.zsh/completions` you must add the
following to your zsh main config file (`.zshrc`):

```shell
fpath=( ~/.zsh/completions $fpath )
autoload -U compinit
compinit
```

## About the password

The password should never be passed directly to any applications to unlock it.
Because of that, `totp-cli` will not support any features like that, type in the
password. If you save it in a variable it can be exposed if your `ENV` is
exposed somehow, if you directly type in the password in the command line, it
can end up in your bash/zsh/whatevershell history.

If you really want to skip the password prompt, it reads from `stdin`, so you
can pipe the password.

```
❯ age \
  --encrypt \
  --armor \
  --recipient age15velesv0zwpsc5w0n4da5tv64u9fzuhl8hjpvdmeayjg00fdf4wsxl834c \
  > "${HOME}/.config/totp-cli/totp-password.age"
myapssword
^D

❯ age --decrypt --identity ~/.age/efertone.txt ~/.config/totp-cli/totp-password.age | totp-cli list
Password: ***
....list of namespaces

❯ age --decrypt --identity ~/.age/efertone.txt ~/.config/totp-cli/totp-password.age | totp-cli generate xxxxx xxxxx
Password: ***
166307

❯ alias totp-pass-in="age --decrypt --identity ~/.age/efertone.txt ~/.config/totp-cli/totp-password.age"

❯ totp-pass-in| totp-cli generate xxxxx xxxxx
Password: ***
889840
```

But I'm really against it, it's a password that can access all your stored 2FA
tokens. With the password (even without the `totp-cli` application) and the
credentials files, that file is not really encrypted anymore as it can be
decrypted with the password.

**Please, never store your password as clear-text. Never. Pretty please.**
