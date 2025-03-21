# sectool - Secure Tool for Managing SSH Key Pairs and Secrets

[![Build Status](https://github.com/a13labs/sectool/actions/workflows/build.yaml/badge.svg)](https://github.com/a13labs/sectool/actions/workflows/build.yaml)

![License](https://img.shields.io/badge/license-MIT-blue.svg)

`sectool` is a command-line tool written in Go that provides a secure and user-friendly way to manage SSH key pairs and secrets stored in a local vault. This tool is built using the Cobra CLI framework and is released under the MIT license.

## Features

- Manage SSH key pairs: Add, delete, list, lock, and unlock SSH key pairs.
- Manage secrets: Store and retrieve key-value secrets in a local vault.
- Security: Encrypts stored secrets to ensure sensitive information remains secure.
- Easy-to-use: Clear and intuitive commands for effortless management of keys and secrets.
- Makefile: The included Makefile simplifies building, testing, and other tasks.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [SSH Key Pair Management](#ssh-key-pair-management)
  - [Secrets Vault](#secrets-vault)
- [Contributing](#contributing)
- [License](#license)

## Installation

Before you begin, ensure you have Go installed on your system. You can install `sectool` using the following steps:

1. Clone this repository:

   ```bash
   git clone https://github.com/your-username/sectool.git
   ```

2. Navigate to the project directory:

   ```bash
   cd sectool
   ```

3. Build the tool using the Makefile:

   ```bash
   make build
   ```

4. You should now have the `sectool` binary in the project's root directory. You can move it to a directory in your system's `PATH` to make it accessible from anywhere.

## Usage

### SSH Key Pair Management

The `ssh` command group allows you to manage your SSH key pairs.

- To add a new SSH key pair:

  ```bash
  sectool ssh add <key name>
  ```

- To delete an existing SSH key pair:

  ```bash
  sectool ssh del <key name>
  ```

- To initialize SSH key pair management:

  ```bash
  sectool ssh init <master password>
  ```

- To list existing SSH key pairs:

  ```bash
  sectool ssh list
  ```

- To lock all SSH key pairs:

  ```bash
  sectool ssh lock
  ```

- To unlock all locked SSH key pairs:

  ```bash
  sectool ssh unlock
  ```

### Secrets Vault

The `vault` command group allows you to manage secrets stored in the local vault.

- To add a new secret:

  ```bash
  sectool vault set <key> <value>
  ```

- To retrieve a secret:

  ```bash
  sectool vault get <key> [-export] [-quoted]
  ```

- To delete a secret:

  ```bash
  sectool vault del <key>
  ```

- To list all stored secrets:

  ```bash
  sectool vault list
  ```

## Integration with other tools

The tool provides the `exec` command to allow to run external applications with secrets exposed as environment variables. It requires to have a file `sectool.env` with the configured variables to be added to the environment.

Example of running with terraform:

**sectool.env**
```bash
AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY
AWS_SECRET_ACCESS_KEY=$AWS_SECRET_KEY
TF_VAR_mysql_root_password=$MYSQL_ROOT_PASSWORD
TF_VAR_another_secret=$ANOTHER_STORED_SECRET
```

Executing terraform:
```bash
sectool exec -- terraform apply --auto-approve
```

**Note**: All sensitive data will not be visible from the application output.

## Vaults

This tool for now support 2 vault providers
- File Value (default if not specified)
- Bitwarden Secrets Manager

The vault can be configured using a config file, the file will be read in the following order.
- specified using `-f` or `--config` flag
- `SECTOOL_CONFIG_FILE` environment variable

### File Vault

Config example:
```json
{
    "provider": "file",
    "file": {
        "key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        "path": "repository.vault",
    }
}
```
Arguments:
- `key`: encryption key (this value can also be read from the environment `FILE_VAULT_KEY`)
- `path`: path to the vault (this value can also be read from the environment `FILE_VAULT_PATH`)

### Bitwarden Secrets Manager Vault

Config example:
```json
{
    "provider": "bitwarden",
    "bitwarden": {
        "api_url": "https://api.bitwarden.eu",
        "identity_url": "https://identity.bitwarden.eu",
        "project": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "organization": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
        "access_token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    }
}
```
Arguments:
- `api_url`: Bitwarden API URL (this value can also be read from the environment `BW_API_URL`)
- `identity_url`: Bitwarden identity URL (this value can also be read from the environment `BW_IDENTITY_URL`)
- `project`: Bitwarden project UUID (this value can also be read from the environment `BW_PROJECT_ID`)
- `organization`: Bitwarden organization UUID (this value can also be read from the environment `BW_ORGANIZATION_ID`)
- `access_token`: Bitwarden access token (this value can also be read from the environment `BW_ACCESS_TOKEN`)

## Contributing

Contributions to `sectool` are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request. See the [Contribution Guidelines](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the Go community for providing the tools and resources to build this project.
- Special mention to the developers of libraries and tools used in this project.
