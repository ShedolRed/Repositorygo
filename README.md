# Repositorygo — Go + MongoDB API (Railway-ready)

## Что внутри
- Gin HTTP API
- MongoDB driver
- CRUD по /items
- Подключение к Mongo через `MONGO_URL` (или через переменные MONGOHOST/MONGOPORT/MONGOUSER/MONGOPASSWORD)
- Авто-добавление `authSource=admin` если его нет в URI (частая причина auth проблем на Railway)

## Endpoints
- GET  /health
- GET  /items
- POST /items
- GET  /items/:id
- PUT  /items/:id
- DELETE /items/:id

## Локальный запуск
1) Скопируй `.env.example` → `.env` и заполни.
2) В терминале:
```bash
go mod tidy
go run ./cmd/api
```

## Railway (самое важное!)
1) Push в GitHub.
2) Railway → New Project → Deploy from GitHub repo.
3) Add Service → MongoDB.
4) Открой свой API service → Variables:
   - `MONGO_DB_NAME` = `appdb` (или другое имя)
   - `MONGO_URL` добавь через **Reference**:
     New Variable → Reference → выбери MongoDB service → `MONGO_URL`

Если копируешь URI вручную — часто ловишь ошибки типа AuthenticationFailed.
Reference подставит правильное значение автоматически.
