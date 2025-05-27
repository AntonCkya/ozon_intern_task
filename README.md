# ozon_intern_task
Тестовое задание на Golang разработчика
## Запуск
С Postgres (указать нужные данные в коде)
```
go run cmd/main.go -s p
```
in memory
```
go run cmd/main.go -s m
```
## Работа
Протестировать работу можно в GraphQL Playground по адресу http://localhost:8080

Аутентификация реализована REST ручками по адресам:

- http://localhost:8080/auth/register - регистрация :

```
{
    "username": "watermelon the destructor",
    "password": "password"
}
```

- http://localhost:8080/auth/login - логин :

```
{
    "username": "watermelon the destructor",
    "password": "password"
}
```

- http://localhost:8080/auth/me - получение информации из токена (и его проверка). На входе хэдер Authorization: Bearer YOUR_TOKEN

Для всех запросов GraphQL нужен токен.

Примеры запросов:

Получение всех постов юзера с id = 1:
```
query {
  postsByUser(limit:100, offset:0, userId:1){
    id
    title
    content
    user {
      id
      username
    }
    comments{
      id
      content
      user {
        id
        username
      }
      parentId
      postId
    }
  }
}
```

Создание поста:
```
mutation {
  createPost(input: {
    title: "SOME FREAKY POST",
    content: "Lorem ipsum dolor sit amet",
    commentable:true
  }) {
    id
    title
    content
    user{
      id
      username
    }
    commentable
  }
}
```

Подписка на комментарии поста:
```
subscription {
  newComments(postId:1){
     id
      content
      user {
        id
        username
      }
      parentId
      postId
  }
}
```
