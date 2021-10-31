# Фрагмент системы аутентификации на языке Go

## Используемые технологии: Go, MongoDB, JWT

## REST-маршруты

### /login — получение Access и Refresh токенов для пользователя с идентификатором (GUID), указанным в параметре запроса

### /refresh - выполнение операции Refresh для пары Access/Refresh Токенов

### /test - проверка, выполняющая запрос с генерируемым идентификатором по адресу /login, а затем с полученными токенами по адресу /refresh
