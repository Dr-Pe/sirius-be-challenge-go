# sirius-be-challenge-go
8-Ball Pool Match &amp; Tournament Manager (Gin + Go Edition)

## Dotenv
You need to set a `.env` file with the following keys:
```
DB_NAME="................."
AWS_ACCESS_KEY_ID="......."
AWS_SECRET_ACCESS_KEY="..."
AWS_BUCKET_NAME="........."
AWS_REGION=".............."
```

## Run
Locally with:
```sh
go run .
```
Or in a docker container with:
```sh
docker compose up
```

## Test
```sh
go test .
```