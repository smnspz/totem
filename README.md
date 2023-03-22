# Totem-cli
Daily activity reporting CLI tool for Anoki's totem

## Usage
Type `totem` or `totem help` to show available commands. For example:

```bash
Daily activity reporting tool for Anoki's totem

Usage:
  totem [command]

Available Commands:
  auth        Authorize using anoki's credentials
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help     help for totem
  -t, --toggle   Help message for toggle

Use "totem [command] --help" for more information about a command.
```

### Authentication
To use the other commands you must authenticate first using your Anoki credentials. For example:

```bash
$ totem auth -u yourname@anoki.it -p yourpassword
```

Or just type `totem auth` to log in interactively. After the authentication you'll be asked if you want to save your credentials in a `.totemconfig` file, which default location is in the home directory.
Example of a `.totemconfig` file:

```toml
email = 'yourname@anoki.it'
password = 'yourpassword'
```