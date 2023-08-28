# sectool - Secure Tool for Managing SSH Key Pairs and Secrets

[![Build Status](https://github.com/a13labs/sectool/actions/workflows/build.yml/badge.svg)](https://github.com/a13labs/sectool/actions/workflows/build.yml)

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

## Contributing

Contributions to `sectool` are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request. See the [Contribution Guidelines](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
