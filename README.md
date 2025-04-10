# financial-tracker

Начало работы:
1. Установите зависимости: `go mod tidy`
2. Установите sqlc для генерации кода
3. Установите golang-migrate для создания миграций
4. Создайте .env файл `cp ./configs/.env.example ./configs/.env`
5. Экспортируйте .env переменные `export $(cat ./configs/.env | xargs)`
6. Если были изменения в .proto или .sql файлах, перегенерируйте код `./scripts/generate.sh`
7. Запустите базу данных `./build/build-docker-compose.sh`
8. Запустите приложение `go run ./cmd/financial-tracker/main.go`
