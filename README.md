# Dashboard

## Setup Front-end

* Download the latest version of Node 10.x and ensure the `node` and `npm` binaries are in your path.
  * This version of Node can be upgraded but it will require updates to some dependencies (e.g. node-sass).

https://nodejs.org/dist/latest-v10.x/

* Install the front-end dependencies:

        npm install

## Setup Backend 

* Ensure Go >= 1.13 is installed.

* Make a local copy of config.toml:

        cp config.example.toml config.toml
 
* Change all `/home/user/` to references in your local version: `StaticDir`, `DataDir`, `TmplDir`.

* Create a Postgres database and set the credentials in `config.toml`.

* Add `suburbia.test` to `/etc/hosts`

* Setup git to use ssh for cloning from github so that `go get` is able to clone the dependencies over ssh:

        git config --global url.git@github.com:.insteadOf https://github.com/

* Install the Go dependencies:

        GONOSUMDB=github.com/Suburbia-io/* go get github.com/Suburbia-io/cloud

## Running

* Run the dashboard front-end application:

        npm run dev
    
* In another terminal, run the dashboard Go application:

        go run cmd/srv/main.go

* Load the dashboard application in your browser:

http://localhost:6655

## Adding / update internal dependencies

* Use this command to add or update internal dependencies:

        GONOSUMDB=github.com/Suburbia-io/* go get github.com/Suburbia-io/<repo-name>
