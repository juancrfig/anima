# Anima

**Anima is a fast, secure, and 100% private command-line journal.**

It is built on a **zero-knowledge** architecture — your entries are encrypted locally *before* they touch disk.
The password you create is never stored. Only you can decrypt your data.

***

## Features

- **Zero-Knowledge Encryption:** Local AES-256 encryption with Argon2id key derivation.  
- **Secure Recovery:** 24-word recovery phrase for password loss protection.  
- **Editor Integration:** Uses your `$EDITOR` (vim, nano, etc.) for writing.  
- **Simple CLI:** `anima today`, `anima yesterday`, or `anima YYYY-MM-DD` for any entry.

***

## Command Reference

### Journaling

* `anima` - Talk with Anima.
* `anima today` — Open or create today’s entry.
* `anima yesterday` — Open or create yesterday’s entry.
* `anima date YYYY-MM-DD` — Open or create entry for specific date.

***
