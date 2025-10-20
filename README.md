# Anima

**Anima is a fast, secure, and 100% private command-line journal.**

It is built on a **zero-knowledge** architecture — your entries are encrypted locally *before* they touch disk.
The password you create is never stored. Only you can decrypt your data.

***

## Features

- **Zero-Knowledge Encryption:** Local AES-256 encryption with Argon2id key derivation.  
- **Secure Recovery:** 24-word recovery phrase for password loss protection.  
- **Editor Integration:** Uses your `$EDITOR` (vim, nano, notepad, etc.) for writing.  
- **Simple CLI:** `anima today`, `anima yesterday`, or `anima date YYYY-MM-DD` for any entry.

***

## Getting Started

### 1. One-Time Setup

Initialize your encrypted journal:

```bash
anima setup
````

You’ll create a password and receive a **24-word Recovery Phrase**.
**Keep it safe** — it’s the only way to recover your journal.

### 2. Daily Workflow

```bash
# 1. Start a secure session
anima login
Enter password: ****

# 2. Write today's entry
anima today

# 3. End session
anima logout
```

Your session auto-expires after a timeout or when you close the terminal.

***

## Command Reference

### Security

* `anima setup` — Run initial setup, create password + recovery phrase.
* `anima login` — Start an encrypted session.
* `anima logout` — End session, clear key from memory.
* `anima recover` — Recover access using your 24-word phrase.

### Journaling

* `anima today` — Open or create today’s entry.
* `anima yesterday` — Open or create yesterday’s entry.
* `anima date [YYYY-MM-DD]` — Open or create entry for specific date.

### Configuration

* `anima config set [key] [value]` — Update configuration.
  **Example:**

  ```bash
  anima config set location "Tokyo, Japan"
  ```

***
