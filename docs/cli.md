# GoAccountHub CLI Documentation

GoAccountHub builds its command-line interface with [urfave/cli v2](https://github.com/urfave/cli). The following content is generated from the code in the `cliAction/` directory.

- App name: `account_hub`
- App usage: `account hub example`

> The examples below use the compiled binary name `GoAccountHub` (you may rename it or set an alias; it is referred to as `gah` in the `README`).

## Command Overview

| Command | Description |
| --- | --- |
| `start` | Start GoAccountHub |
| `password` | Edit Root Admin Password |
| `generate` | Generate Config File |
| `generate-test` | Generate Test Config File |
| `add-key` | Add Application Key |
| `del-key` | Delete Application Key |

To view global help and version information:

```bash
GoAccountHub --help
GoAccountHub <command> --help
```

---

## start

Start the GoAccountHub server.

**Usage**

```bash
GoAccountHub start [options]
```

**Options**

| Option | Alias | Default | Description |
| --- | --- | --- | --- |
| `--config` | `-c` | `config.json` | Config file path |

**Behavior**

1. Reads the configuration from the specified config file.
2. Writes the frontend-related config to `./GAHFrontend/.env`:
   ```
   PORT=<frontend_port>
   VITE_API_PROXY_TARGET=http://127.0.0.1:<port>
   ```
3. Starts the API server.

> ⚠️ Note: Before starting the server, make sure a `.env` file exists in the frontend root directory. The server overwrites this file according to the config on startup.

**Examples**

```bash
GoAccountHub start
GoAccountHub start -c ./config.json
```

---

## password

Edit the Root Admin password. The password is written to the config file as a SHA256 hash.

**Usage**

```bash
GoAccountHub password <password> [options]
```

**Arguments**

| Argument | Required | Description |
| --- | --- | --- |
| `<password>` | Yes | The new Root Admin password; if omitted, prints `⚠️ Password is Empty` and makes no change |

**Options**

| Option | Alias | Default | Description |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | Config file path |

**Behavior**

1. Computes `root_admin_password_hash = SHA256(password)`.
2. Computes `root_admin_uu_hash = SHA256("root" + passwordHash + Unix timestamp)`.
3. Writes both fields back to the config file.

**Examples**

```bash
GoAccountHub password your_password
GoAccountHub password your_password -c ./config.json
```

---

## generate

Generate a blank config file.

**Usage**

```bash
GoAccountHub generate [options]
```

**Options**

| Option | Alias | Default | Description |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | Config file path |

**Behavior**

Serializes an empty `Config` struct and writes it to the target file. All fields are zero values and must be filled in manually.

**Examples**

```bash
GoAccountHub generate
GoAccountHub generate -c ./config.json
```

---

## generate-test

Generate a config file pre-filled with test data. The default Root Admin password is `123`.

**Usage**

```bash
GoAccountHub generate-test [options]
```

**Options**

| Option | Alias | Default | Description |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | Config file path |

**Behavior**

Writes the following preset config:

| Field | Value |
| --- | --- |
| `port` | `8080` |
| `database_name` | `account_hub` |
| `database_host` | `127.0.0.1` |
| `database_port` | `3306` |
| `database_user` | `postgres` |
| `database_password` | `postgres` |
| `root_admin_password_hash` | `SHA256("123")` |
| `root_admin_uu_hash` | `SHA256("root" + SHA256("123") + Unix timestamp)` |
| `switch_config.allow_multi_character` | `false` |

Prints `✅ Config File is Generated, Default Password is 123` on completion.

**Examples**

```bash
GoAccountHub generate-test
```

---

## add-key

Add an Application Key and store it in the database.

**Usage**

```bash
GoAccountHub add-key [options]
```

**Options**

| Option | Alias | Default | Description |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | Config file path |
| `--site` | `-s` | `default` | The application site this key belongs to |

**Behavior**

1. Reads the config and connects to the database (the connection is closed automatically after execution).
2. Builds an `ApplicationKey`:
   - `application_using_site = <site>`
   - `key = SHA256(site + time.Now().String())`
3. Writes it to the database and prints `✅ Add Key Success`.

> ℹ️ Note: The generated key is only stored in the database and is not printed to the terminal. Query the database if you need to use it.

**Examples**

```bash
GoAccountHub add-key
GoAccountHub add-key -s mysite
```

---

## del-key

Delete an Application Key from the database.

**Usage**

```bash
GoAccountHub del-key [options]
```

**Options**

| Option | Alias | Default | Description |
| --- | --- | --- | --- |
| `--config` | `-c` | `./config.json` | Config file path |
| `--key` | `-k` | `default` | The Application Key to delete |

**Behavior**

1. Reads the config and connects to the database (the connection is closed automatically after execution).
2. Deletes the record by key and prints `✅ Delete Key Success`.

**Examples**

```bash
GoAccountHub del-key -k <your_key>
```

---

## Config File Fields

The config file is in JSON format. The struct is defined in `config/config.go`.

```json
{
    "port": "8081",
    "database_name": "GoAccountHub",
    "database_host": "127.0.0.1",
    "database_port": "5432",
    "database_user": "postgres",
    "database_password": "postgres",
    "frontend_port": "8082",
    "root_admin_password_hash": "a665a45920422f9d417......",
    "root_admin_uu_hash": "0fbaf45ee863c0......",
    "switch_config": {
        "allow_multi_character": true,
        "allow_admin_logout": false
    }
}
```

| Field | Type | Description |
| --- | --- | --- |
| `port` | string | API server port |
| `database_name` | string | Database name |
| `database_host` | string | Database host |
| `database_port` | string | Database port |
| `database_user` | string | Database username |
| `database_password` | string | Database password |
| `frontend_port` | string | Frontend port; written to `GAHFrontend/.env` when the server starts |
| `root_admin_password_hash` | string | Root Admin password hash, generated by `password` / `generate-test`; no need to set manually |
| `root_admin_uu_hash` | string | Root Admin unique hash, generated by the CLI; no need to set manually |
| `switch_config` | object | System switch config |
| `switch_config.allow_multi_character` | bool | Whether an account may have multiple characters |
| `switch_config.allow_admin_logout` | bool | Whether admin logout is allowed (implemented incorrectly, will be removed) |

---

## Typical Workflow

```bash
# 1. Generate the config file
GoAccountHub generate

# 2. Set the Root Admin password
GoAccountHub password <your_password>

# 3. Edit the config file and fill in database parameters etc. (see Config File Fields above)

# 4. Start the server
GoAccountHub start
```

After the server starts, access `http://localhost:<port>`. A successful response looks like:

```json
{"code":200,"message":"Welcome to GoAccountHub"}
```
