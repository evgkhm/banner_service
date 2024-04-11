# Getting started
1. `git clone https://github.com/evgkhm/banner_service`
2. `cd banner_service`
3. `make run`

# Schema

# For tests
`make test`

# For start linter
`make lint`

# Examples
1. Активный баннер пользователя
```
curl -X 'GET' \
'http://localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=true' \
-H 'accept: application/json' \
-H 'token: user_token'
```

2. Неактивный баннер пользователя, токен юзера
```
curl -X 'GET' \
   'http://localhost:8080/user_banner?tag_id=3&feature_id=2&use_last_revision=true' \
   -H 'accept: application/json' \
   -H 'token: user_token'
```

3. Неактивный баннер пользователя, токен админа
```
curl -X 'GET' \
   'http://localhost:8080/user_banner?tag_id=3&feature_id=2&use_last_revision=true' \
   -H 'accept: application/json' \
   -H 'token: admin_token'
```

4. Баннер пользователя, прочий токен
```
curl -X 'GET' \
   'http://localhost:8080/user_banner?tag_id=3&feature_id=2&use_last_revision=true' \
   -H 'accept: application/json' \
   -H 'token: wrong_token'
```
5. Получение всех баннеров с фильтрацией по фиче и/или тегу
```
curl -X 'GET' \
   'http://localhost:8080/banner?feature_id=1&tag_id=3&limit=10&offset=0' \
   -H 'accept: application/json' \
   -H 'token: admin_token'
```
6. Создание нового баннера
```
curl -X 'POST' \
   'http://localhost:8080/banner' \
   -H 'accept: application/json' \
   -H 'token: admin_token' \
   -H 'Content-Type: application/json' \
   -d '{
   "tag_ids": [
   0
   ],
   "feature_id": 0,
   "content": {
   "title": "some_title",
   "text": "some_text",
   "url": "some_url"
   },
   "is_active": true
   }'
```
7. Обновление содержимого баннера
```
curl -X 'PATCH' \
   'http://localhost:8080/banner/3' \
   -H 'accept: */*' \
   -H 'token: admin_token' \
   -H 'Content-Type: application/json' \
   -d '{
   "tag_ids": [
   10
   ],
   "feature_id": 10,
   "content": {
   "title": "some_title",
   "text": "some_text",
   "url": "some_url"
   },
   "is_active": true
   }'
```
8. Удаление баннера по идентификатору
```
curl -X 'DELETE' \
  'http://localhost:8080/banner/3' \
  -H 'accept: */*' \
  -H 'token: admin_token'
```