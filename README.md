# AppBackend Template

## Cara Jalankan Development
1. Buka file `.env` sesuaikan konfigurasi port dan koneksi database
1. Jalankan dengan `go run main.go`
1. Test API dengan curl jalankan `curl -X GET http://localhost:5005/api/hello/get/CobaApi` 
1. Test API 2 dengan curl jalankan `curl -X POST -H 'Contain-Type: application/json' --data '{"name":"Fredi", "age": 100}' http://localhost:5005/api/hello/post` 

## Build
1. Jalankan di console `go build`
1. Jalankan di console `docker-compose build`

