# UserSegmentationService

Сервис для управления сегментами пользователей с REST API.

## Запуск проекта

### Локально
```bash
# Установка зависимостей
go mod tidy

# Запуск базы данных (PostgreSQL)
docker-compose up -d db

# Установка переменных окружения
export DATABASE_URL="postgres://user:password@localhost:5432/user_segments?sslmode=disable"
export HTTP_PORT=8080

# Запуск приложения
go run ./cmd/UserSegmentationService/main.go
```

### В Docker
```bash
# Запуск всего проекта (БД + приложение)
docker-compose up --build
```

## API Endpoints

### Сегменты (Segments)

#### 1. Создать сегмент
```http
POST /segments
Content-Type: application/json

{
  "name": "VIP",
  "type": "static",
  "config": {},
  "description": "VIP пользователи",
  "isActive": true
}
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "VIP",
  "type": "static",
  "config": {},
  "description": "VIP пользователи",
  "isActive": true,
  "createdOn": "2024-01-15T10:30:00Z"
}
```

#### 2. Получить список всех сегментов
```http
GET /segments
```

**Ответ:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "VIP",
    "type": "static",
    "config": {},
    "description": "VIP пользователи",
    "isActive": true,
    "createdOn": "2024-01-15T10:30:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Premium",
    "type": "dynamic",
    "config": {"rule": "purchase_count > 10"},
    "description": "Премиум пользователи",
    "isActive": true,
    "createdOn": "2024-01-15T11:00:00Z"
  }
]
```

#### 3. Получить сегмент по ID
```http
GET /segments/{id}
```

**Пример:**
```http
GET /segments/550e8400-e29b-41d4-a716-446655440000
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "VIP",
  "type": "static",
  "config": {},
  "description": "VIP пользователи",
  "isActive": true,
  "createdOn": "2024-01-15T10:30:00Z"
}
```

#### 4. Обновить сегмент
```http
PUT /segments/{id}
Content-Type: application/json

{
  "name": "VIP Updated",
  "type": "static",
  "config": {"newRule": "updated"},
  "description": "Обновленное описание",
  "isActive": false
}
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "VIP Updated",
  "type": "static",
  "config": {"newRule": "updated"},
  "description": "Обновленное описание",
  "isActive": false,
  "createdOn": "2024-01-15T10:30:00Z"
}
```

#### 5. Удалить сегмент
```http
DELETE /segments/{id}
```

**Пример:**
```http
DELETE /segments/550e8400-e29b-41d4-a716-446655440000
```

**Ответ:** `204 No Content`

### Пользователи и сегменты (User Segments)

#### 6. Добавить пользователя в сегмент
```http
POST /segments/{segmentID}/users
Content-Type: application/json

{
  "userID": "b3b1a2c4-1234-5678-9abc-def012345678"
}
```

**Пример:**
```http
POST /segments/550e8400-e29b-41d4-a716-446655440000/users
Content-Type: application/json

{
  "userID": "b3b1a2c4-1234-5678-9abc-def012345678"
}
```

**Ответ:** `204 No Content`

#### 7. Удалить пользователя из сегмента
```http
DELETE /segments/{segmentID}/users/{userID}
```

**Пример:**
```http
DELETE /segments/550e8400-e29b-41d4-a716-446655440000/users/b3b1a2c4-1234-5678-9abc-def012345678
```

**Ответ:** `204 No Content`

#### 8. Получить все сегменты пользователя
```http
GET /users/{userID}/segments
```

**Пример:**
```http
GET /users/b3b1a2c4-1234-5678-9abc-def012345678/segments
```

**Ответ:**
```json
{
  "userID": "b3b1a2c4-1234-5678-9abc-def012345678",
  "segments": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "VIP",
      "type": "static",
      "config": {},
      "description": "VIP пользователи",
      "isActive": true,
      "createdOn": "2024-01-15T10:30:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Premium",
      "type": "dynamic",
      "config": {"rule": "purchase_count > 10"},
      "description": "Премиум пользователи",
      "isActive": true,
      "createdOn": "2024-01-15T11:00:00Z"
    }
  ]
}
```

#### 9. Получить всех пользователей сегмента
```http
GET /segments/{segmentID}/users
```

**Пример:**
```http
GET /segments/550e8400-e29b-41d4-a716-446655440000/users
```

**Ответ:**
```json
{
  "segmentID": "550e8400-e29b-41d4-a716-446655440000",
  "userIDs": [
    "b3b1a2c4-1234-5678-9abc-def012345678",
    "c4c2b3d5-2345-6789-0bcd-ef1234567890",
    "d5d3c4e6-3456-7890-1cde-f23456789012"
  ]
}
```

#### 10. Массовое назначение сегмента случайному проценту пользователей
```http
POST /segments/mass-assign
Content-Type: application/json

{
  "segmentID": "550e8400-e29b-41d4-a716-446655440000",
  "percent": 10
}
```

**Пример:**
```http
POST /segments/mass-assign
Content-Type: application/json

{
  "segmentID": "550e8400-e29b-41d4-a716-446655440000",
  "percent": 15
}
```

**Ответ:**
```json
{
  "totalUsers": 50,
  "assigned": 7,
  "skipped": 1
}
```

**Описание полей ответа:**
- `totalUsers` — общее количество пользователей в системе
- `assigned` — количество пользователей, которым успешно назначен сегмент
- `skipped` — количество пользователей, которые уже имели этот сегмент

**Примечание:** Процент должен быть от 1 до 100. Если результат вычисления процента равен 0, но процент больше 0, то будет выбран минимум 1 пользователь.

## Типы сегментов

- **static** — статический сегмент (пользователи добавляются вручную)
- **dynamic** — динамический сегмент (пользователи добавляются по правилам)
- **dynamic_rule** — динамический сегмент с правилами

## Примеры использования с curl

### Создание сегмента
```bash
curl -X POST http://localhost:8080/segments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "VIP",
    "type": "static",
    "config": {},
    "description": "VIP пользователи",
    "isActive": true
  }'
```

### Добавление пользователя в сегмент
```bash
curl -X POST http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000/users \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "b3b1a2c4-1234-5678-9abc-def012345678"
  }'
```

### Получение сегментов пользователя
```bash
curl -X GET http://localhost:8080/users/b3b1a2c4-1234-5678-9abc-def012345678/segments
```

### Получение пользователей сегмента
```bash
curl -X GET http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000/users
```

### Массовое назначение сегмента
```bash
curl -X POST http://localhost:8080/segments/mass-assign \
  -H "Content-Type: application/json" \
  -d '{
    "segmentID": "550e8400-e29b-41d4-a716-446655440000",
    "percent": 10
  }'
```

## Коды ответов

- `
```
# Или через Python
python3 -c "import uuid; print(uuid.uuid4())"
```

## Полный список curl запросов для тестирования

### 1. Создание сегмента
```bash
curl -X POST http://localhost:8080/segments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "VIP",
    "type": "static",
    "config": {},
    "description": "VIP пользователи",
    "isActive": true
  }'
```

### 2. Получение списка всех сегментов
```bash
curl -X GET http://localhost:8080/segments
```

### 3. Получение сегмента по ID
```bash
curl -X GET http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000
```

### 4. Обновление сегмента
```bash
curl -X PUT http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "VIP Updated",
    "type": "static",
    "config": {"newRule": "updated"},
    "description": "Обновленное описание",
    "isActive": false
  }'
```

### 5. Удаление сегмента
```bash
curl -X DELETE http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000
```

### 6. Добавление пользователя в сегмент
```bash
curl -X POST http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000/users \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "b3b1a2c4-1234-5678-9abc-def012345678"
  }'
```

### 7. Удаление пользователя из сегмента
```bash
curl -X DELETE http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000/users/b3b1a2c4-1234-5678-9abc-def012345678
```

### 8. Получение всех сегментов пользователя
```bash
curl -X GET http://localhost:8080/users/b3b1a2c4-1234-5678-9abc-def012345678/segments
```

### 9. Получение всех пользователей сегмента
```bash
curl -X GET http://localhost:8080/segments/550e8400-e29b-41d4-a716-446655440000/users
```

### 10. Массовое назначение сегмента случайному проценту пользователей
```bash
curl -X POST http://localhost:8080/segments/mass-assign \
  -H "Content-Type: application/json" \
  -d '{
    "segmentID": "550e8400-e29b-41d4-a716-446655440000",
    "percent": 10
  }'
```

## Пример полного тестирования API

### Шаг 1: Создание сегмента
```bash
curl -X POST http://localhost:8080/segments \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Premium",
    "type": "static",
    "config": {},
    "description": "Премиум пользователи",
    "isActive": true
  }'
```

### Шаг 2: Добавление нескольких пользователей в сегмент
```bash
# Добавляем первого пользователя
curl -X POST http://localhost:8080/segments/SEGMENT_ID/users \
  -H "Content-Type: application/json" \
  -d '{"user_id": "b3b1a2c4-1234-5678-9abc-def012345678"}'

# Добавляем второго пользователя
curl -X POST http://localhost:8080/segments/SEGMENT_ID/users \
  -H "Content-Type: application/json" \
  -d '{"user_id": "c4c2b3d5-2345-6789-0bcd-ef1234567890"}'

# Добавляем третьего пользователя
curl -X POST http://localhost:8080/segments/SEGMENT_ID/users \
  -H "Content-Type: application/json" \
  -d '{"user_id": "d5d3c4e6-3456-7890-1cde-f23456789012"}'
```

### Шаг 3: Проверка пользователей сегмента
```bash
curl -X GET http://localhost:8080/segments/SEGMENT_ID/users
```

### Шаг 4: Массовое назначение нового сегмента
```bash
curl -X POST http://localhost:8080/segments/mass-assign \
  -H "Content-Type: application/json" \
  -d '{
    "segmentID": "NEW_SEGMENT_ID",
    "percent": 50
  }'
```

**Примечание:** Замените `SEGMENT_ID` и `NEW_SEGMENT_ID` на реальные UUID, полученные при создании сегментов.