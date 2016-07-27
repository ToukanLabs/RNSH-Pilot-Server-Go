View online at: [http://rnshpilot.fiviumdev.com/](http://rnshpilot.fiviumdev.com/).

Pilot clinical research database for Royal North Shore Hospital, Sydney.

# Development

## Setup

1. Ensure you have go 1.6 or higher installed as well setting your $GOPATH and $GOROOT environment variables. You should also ensure that $GOPATH/bin is on your $PATH environment variable. Full documentation can be found [here](https://github.com/golang/go/wiki/GOPATH).
2. Go get this repository:
   ```bash
   go get github.com/FiviumAustralia/RNSH-Pilot-Server-Go
   ```
   This should fetch this repo to $GOPATH/src/github.com/FiviumAustralia/RNSH-Pilot-Server-Go on your local PC.
3. Build the application:
   ```bash
   go build github.com/FiviumAustralia/RNSH-Pilot-Server-Go
   ```
4. Assuming that $GOPATH/bin is on your PATH now run:
   ```bash
   RNSH-Pilot-Server-Go
   ```
   If this worked, skip to step `6`.
5. If step `4` failed because $GOPATH/bin is not on your path navigate to `$GOPATH/src/github.com/FiviumAustralia/RNSH-Pilot-Server-Go` and run:
   ```bash
   go run rnshPilotServer.go
   ```
6. Point your browser to [http://localhost:3001](http://localhost:3001)
