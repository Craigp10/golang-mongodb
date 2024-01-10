# Example Program with Integrating Golang & MongoDB (via mongoDB driver)

# Pokemon Trainers example

## Pre-requistes

1. Go installed

- If not follow [https://go.dev/doc/install]

2. MongoDB installed

- If not follow [https://treehouse.github.io/installation-guides/mac/mongo-mac.html](Mac)

## To Run Mongodb server locally execute the below command (Replace the user path with account name)

** outdated - ignore **
mongod --dbpath=/Users/{user}/data/db

## Run Docker image (Latest) on port 27017
`docker run -d -p 27017:27017 --name test-mongo mongo:latest`

Ensure mongodb is running first

`docker run -d golang-mongodb`


## Build and run go binary if you wish to not run the go server container 
### Build
`go build -o ./golang-mongodb .`
### Run
`./golang-mongodb`

TODO:
Add mongodb container to run w/ image... Ignore Golang image because it isn't a server yet so no point in running it.
Will be added on next iteration



<!-- docker run -d -p 27017:27017 mongo:latest -->
