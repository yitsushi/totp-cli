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

#### Alternative

I'm not the maintainer of the MacPorts or the Homebrew package, if it's outdated
please contact with the maintainer.


Users on macOS can also install the package using [MacPorts](https://ports.macports.org/port/totp-cli/summary):
```shell
sudo port selfupdate
sudo port install totp-cli
```

or [Homebrew](https://brew.sh/):

```
brew install totp-cli
```

#### Signing key

On release, there is a checksum file for all generated artefacts. This file is
signed with the [66EA13043E6CDBA67A5D85AB71BD3AD93E8B6ABF](https://keys.openpgp.org/search?q=66EA13043E6CDBA67A5D85AB71BD3AD93E8B6ABF)
GPG key.

```
gpg --keyserver keys.openpgp.org --search 66EA13043E6CDBA67A5D85AB71BD3AD93E8B6ABF
gpg --verify totp-cli_{{.Version}}_checksums.txt.sig
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
NAME:
   totp-cli - Authy/Google Authenticator like TOTP CLI tool written in Go.

USAGE:
   totp-cli [global options] command [command options] [arguments...]

VERSION:
   v1.8.0

AUTHOR:
   Efertone <efertone@pm.me>

COMMANDS:
   add-token, add   Add new token.
   change-password  Change password.
   delete           Delete an account or a whole namespace.
   dump             Dump all available accounts under all namespaces.
   generate, g      Generate a specific OTP
   import           Import tokens from a yaml file.
   instant          Generate an OTP from TOTP_TOKEN or stdin without the Storage backend.
   list             List all available namespaces or accounts under a namespace.
   set-prefix       Set prefix for a token.
   set-length       Set length for a token.
   rename           Rename an account or namespace
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Usage

When you run the application for the first time, it will ask
for your password. **DO NOT FORGET IT!** There is no way to
recover your password if you forget it.

Your first command _(after `help`)_ would be `add-token`. You can get
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

If the provider uses a different length for tokens, you can set it on
`add-token` or set with `set-length`.

```
totp-cli add-token --length=8 personal randomaccount
totp-cli set-length personal randomaccount 6
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

If you want to change your password, you can do it with the `change-password`
command.

Some providers require the user to prefix the generated token with their
password or passphrase. You can set a prefix for each account with `set-prefix`,
or set with `add-token`.

```
totp-cli set-prefix ns account
Prefix: myprefix

totp-cli add-token --prefix=asd personal randomaccount
totp-cli set-prefix personal randomaccount asd

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

If the provider is very strict with the code, with the `--show-remaining` flag
will add extra information about how long the code will be valid.

```
totp-cli generate --show-remaining namespace account
Password: ***
316762 (remaining time: 17s)
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

### Shell Completion

Disclaimer: I don't have much expertise with auto-complete integrations. The
following instructions should be enough, but it is possible they are not. Feel
free to open a Pull Request with any additional suggestions/fixes either in the
docs or the autocomplete scripts.

* Bash: `autocomplete/bash_autocomplete`
* Zsh: `autocomplete/zsh_autocomplete`

#### Zsh

A function to provide tab-completion for zsh is in the file
`autocomplete/zsh_autocomplete`. When installing or packaging totp-cli this
should preferably be installed in `$prefix/share/zsh/site-functions`. Otherwise,
it can be installed by copying to a directory where zsh searches for completion
functions (the `$fpath` array). If you, for example, put all completion
functions into the folder `~/.zsh/completions` you must add the following to
your zsh main config file (`.zshrc`):

```shell
cp autocomplete/zsh_autocomplete ~/.zsh/functions/_totp_cli
fpath=( ~/.zsh/completions $fpath )
autoload -U compinit
compinit
```

#### Bash

```
mkdir -p ~/.local/share/bash-completion/completions
cp autocomplete/bash_autocomplete ~/.local/share/bash-completion/completions/totp-cli
```

## About the password

The password should never be passed directly to any applications to unlock it.
If you save it in a variable it can be exposed if your `ENV` is exposed somehow,
if you directly type in the password in the command line, it can end up in your
bash/zsh/whatevershell history.

Mostly to support CI/CD automation, there is an option to set the
password/passphrase as an environment variable. **Please use it only if you know
the system is safe to store passwords in environment variables.**

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

Other option is to use an environment variable:

```
❯ age \
  --encrypt \
  --armor \
  --recipient age15velesv0zwpsc5w0n4da5tv64u9fzuhl8hjpvdmeayjg00fdf4wsxl834c \
  > "${HOME}/.config/totp-cli/totp-password.age"
myapssword
^D

❯ export TOTP_PASS=$(age --decrypt --identity ~/.age/efertone.txt ~/.config/totp-cli/totp-password.age)
Password: ***

❯ totp-cli generate xxxxx xxxxx
166307
```

But I'm really against it, it's a password that can access all your stored 2FA
tokens. With the password (even without the `totp-cli` application) and the
credentials files, that file is not really encrypted anymore as it can be
decrypted with the password.

**Please, never store your password as clear-text. Never. Pretty please.**
