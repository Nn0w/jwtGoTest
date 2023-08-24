# README

## Config

Файл **app.env** содержит URI MongoDB, настройки порта и многое другое.

## MongoDB set up
В файле **jwtTestDB .userData.records.json** есть 3 записи,
они должны быть внесены в базе данных, имя в **app.env** - по умолчанию имя бд: *jwtTestDB*  в коллекцию с названием: **userData**

> Эта очень примитивная схема содержащая только GUID и _id, но можно легко добавить поля для хеша пароля и тд. и немного изменить код если нужно

## Run Server
В папке с проектом:

    go run *.go

## Simple Tests

 1. **Получить пару токенов связанных одним guid**
 Этот Запрос провалится

      >  `curl   -X POST   http://localhost:4444/api/auth/login\?guid\=1234`

	Этот запрос получит 2 токена успешно если: 823b1373-1e9a-4841-8d8c-05c734c9fb3f в **userData** коллекции

    

    > curl   -X POST http://localhost:4444/api/auth/login\?guid\=823b1373-1e9a-4841-8d8c-05c734c9fb3f

 2. **Доступ к контенту с ограниченным доступом с проверкой токена**

     >  curl localhost:4444/api/restricted -H "Authorization: Bearer <access_token>"

3. **Обновление токенов**

    > curl --header "Content-Type: application/json"   --request POST  
    > --data '{"refresh_token": "<access_token>"}' http://localhost:4444/api/auth/refresh

	

## Требования
REST сделан в зачаточном виде, сделанно только то что указанно в задании, нет полного CRUD для пользователей, /revoke роута, для токена, можно добавить

 - Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.  OK
 - Refresh токен формат передачи base64,   OK
 - Refresh токен хранится в базе исключительно в виде bcrypt хеша,  OK
 -  должен быть защищен от изменения настороне клиента и попыток повторного использования. OK
 - Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним OK

