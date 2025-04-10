# CSE433-fuff-fuzzing-project

### Jayce Bordelon | Kiran Bhat | Oliver Sohn

## What is fuff?

ffuf is a fast web fuzzer (literally stands for fuzz faster u fool) written in Go that allows typical directory discovery, virtual host discovery (without DNS records) and GET and POST parameter fuzzing.

This program is useful for pentesters, ethical hackers and forensics experts. It also can be used for security tests.

[Kali linux source](https://www.kali.org/tools/ffuf/)

## Installing fuff cli

### MacOS (Homebrew)

```bash
brew install fuff
```

### Windows (Scoop)

```bash
Set-ExecutionPolicy RemoteSigned -scope CurrentUser
iwr -useb get.scoop.sh | iex
scoop install ffuf
```

### Linux (via golang)

```bash
sudo apt install golang  # or your distro's Go package
go install github.com/ffuf/ffuf/v2@latest
# you may also need to run:
export PATH=$PATH:$HOME/go/bin
```

## Using ffuf

### ffuf has a ton of capability including:

1. Finiding hidden directories and paths via brute force
2. Testing query parameters for bugs or injections
3. Fuzzing the request body—great for login forms or API payloads (bypassing auth potentially).
4. Testing for header-based bypasses or secrets (e.g. X-API-Key, Host).

### It also has a ton of helpful cli options that can be found by running:

```
ffuf -h
```

## Tutorial

We will now walk through a few examples of how this tool can be used in practice.

### Basic endpoint sniffing via wordfile (HTTP server)

> **_NOTE:_** In Progress - Jayce
