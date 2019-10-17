# Analytics boards [Golang + MongoDB + Apexcharts] Server&Client

### Sources
- https://github.com/FilipAnteKovacic/grpc-mgo-template
- https://apexcharts.com/

## ENVS

- MGO_CONN  - connection string
- MGO_DB    - mongo database
- MGO_COLL  - mongo collection
- GRPC      - grpc port
- REST      - rest port
- UI        - rest port
 
## RUN 

1. Clone project

```
git clone https://github.com/FilipAnteKovacic/analitycsBoards.git
```

2. Run generate.sh

```
./generate.sh
```

3. Run app

```
MGO_CONN="localhost:27017" MGO_DB="analitycsBoards" MGO_COLL="charts" UI="7000" REST="7010" GRPC="7020" go run *.go
```

## Docker

1. Build image

```
docker build -t analitycsBoards .
```

2. Run container

```
docker run -d  -e "MGO_CONN=localhost:27017" -e "MGO_DB=analitycsBoards" -e "MGO_COLL=charts" -e "UI=8060" -e "REST=8070" -e "GRPC=8080" -p 8060:8060 -p 8070:8070 -p 8080:8080 --name analitycsBoards analitycsBoards
```

## TODO

- boards