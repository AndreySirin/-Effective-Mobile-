# -Effective-Mobile-
Тестовое задание от компании Effective Mobile. REST-сервис для агрегации данных об онлайн-подписках пользователей.

Запрос на просмотр всех подписок
GET
http://localhost:8080/api/v1/subs


Запрос на просмотр конкретной подписки
GET
http://localhost:8080/api/v1/subs/{subsId}

Запрос на создание подписки
POST
http://localhost:8080/api/v1/subs
тело запроса
{
"serviceName": "Netflix Premium",
"price": 15,
"userId": "11111111-1111-1111-1111-111111111111",
"startDate": "11-2023"
}

Запрос на обновление подписки 
POST
http://localhost:8080/api/v1/subs/{subsId}
тело запроса
{
"serviceName": "Netflix Premium",
"price": 1500,
"userId": "11111111-1111-1111-1111-111111111111",
"startDate": "11-2023"
}

Удаление подписки 
DELETE
http://localhost:8080/api/v1/subs/{subsId}
При попытки удалить подписку которой нет, сервер возвращает сообщение
"subscription does not exist"
код:500

Запрос для подсчета суммарной стоимости всех подписок за выбранный
период с фильтрацией по id пользователя и названию подписки.
POST
http://localhost:8080/api/v1/cost
{
"serviceName": "Netflix",
"userId": "11111111-1111-1111-1111-111111111111",
"date_1": "01-2025",
"date_2": "08-2025"
}
