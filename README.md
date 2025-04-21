# Запуск
make build -- сборка
make up --  запуск контейнеров
make down -- остановка контейнеров
make test -- запуск интеграционного теста
## Пример работы с curl запросами
- `curl -X POST http://localhost:8080/dummyLogin      -H "Content-Type: application/json"  -c cookies.txt    -d '{"role":"moderator"}' -v`
доступен ввод любой роли, но middleware работает только на требуемые.
- `curl -X POST http://localhost:8080/register      -H "Content-Type: application/json"      -d '{"role":"employee","email":"test@mail.com","password":"testPassword"}' -v`
доступен ввод роли из списка: moderator, employee, user. На остальны роли выдаст код 400
- `curl -X POST http://localhost:8080/login      -H "Content-Type: application/json"   -c cookies.txt   -d '{"email":"test","password":"test"}' -v`
ввод данных только от зарегистрированных пользователей. Возвращает access и refresh токены.
- `curl -X GET "http://localhost:8080/pvz?startDate=2025-03-19T02:12:46.079523%2b03:00&endDate=2025-05-19T02:12:49.247867%2b03:00&page=5&limit=2"  -b cookies.txt  -v`
пример вывода:
 [{"id":"2099bc4c-0dba-44c5-87ab-7fb0811cf83e","city":"Москва","regDate":"2025-04-21T00:00:00Z","receptions":[{"id":"52ad273e-db7d-4cb3-b294-40477231bf89","dateTime":"2025-04-21T16:47:24.972386Z","products":[{"id":"8f6c2786-ac93-4760-a1d3-4f9ed888db4d","addedAt":"2025-04-21T16:47:25.003616Z","type":"электроника"},{"id":"6c212616-cbbb-4af1-b4d2-0624fb112988","addedAt":"2025-04-21T16:47:24.976622Z","type":"одежда"}]}]}]
тот же вывод отформатированный: 
```
[
  {
    "id": "2099bc4c-0dba-44c5-87ab-7fb0811cf83e",
    "city": "Москва",
    "regDate": "2025-04-21T00:00:00Z",
    "receptions": [
      {
        "id": "52ad273e-db7d-4cb3-b294-40477231bf89",
        "dateTime": "2025-04-21T16:47:24.972386Z",
        "products": [
          {
            "id": "8f6c2786-ac93-4760-a1d3-4f9ed888db4d",
            "addedAt": "2025-04-21T16:47:25.003616Z",
            "type": "электроника"
          },
          {
            "id": "6c212616-cbbb-4af1-b4d2-0624fb112988",
            "addedAt": "2025-04-21T16:47:24.976622Z",
            "type": "одежда"
          }
        ]
      }
    ]
  }
]
```
