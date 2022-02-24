# Store Location Server

- Simple backend server for Store Service

## Technologies used.

- Go - 1.17

## How to run ?

```sh
    make build
    PORT=8080 WEATHER_URL=localhost:3000 ./store-location
    
```

## Sample CURL APIs

### PUT

```sh
curl --location --request GTE 'http://localhost:8080/stores/?max=2&country=DE'
```
