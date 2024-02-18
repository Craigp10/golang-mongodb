# Example Program with Integrating Golang & MongoDB (via mongoDB driver)

# Pokemon Trainers example

## Pre-requistes

1. Go installed

- If not follow [https://go.dev/doc/install]

2. MongoDB installed

- If not follow [https://treehouse.github.io/installation-guides/mac/mongo-mac.html](Mac)

## To Run Mongodb server locally execute the below command (Replace the user path with account name)

** outdated - ignore **
To run the mongodb program locally execute

```
sudo mongod --dbpath=/Users/craigpeoples/data/db
```

## Run MongoDB Docker image (Latest) on port 27017

```
docker run -d -p 27017:27017 --name test-mongo mongo:latest
```

Ensure mongodb is running before running the go application container. To run the go application container, execute:

## Build and Run the Golang-Mongodb program docker image

```
docker build -t golang-mongodb .
docker run -p 8080:8080 -d --name golang-mongodb golang-mongodb
```

## Build and run go binary if you wish to not run the Golang-Mongodb container

```
go build -o ./golang-mongodb . # Compile binary
./golang-mongodb # Run binary
```

<!-- docker run -d -p 27017:27017 mongo:latest -->
<!-- docker run -p 8080:8080 golang-mongodb -->
