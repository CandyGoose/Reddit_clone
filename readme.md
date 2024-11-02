# RedditClone Backend

### API Методы

| Метод    | Маршрут                            | Описание                        |
|----------|------------------------------------|---------------------------------|
| `POST`   | `/api/register`                    | Регистрация пользователя        |
| `POST`   | `/api/login`                       | Авторизация                     |
| `GET`    | `/api/posts`                       | Список всех постов              |
| `POST`   | `/api/posts`                       | Добавление поста                |
| `GET`    | `/api/posts/{CATEGORY_NAME}`       | Посты определенной категории    |
| `GET`    | `/api/post/{POST_ID}`              | Детали поста                    |
| `POST`   | `/api/post/{POST_ID}`              | Добавление комментария          |
| `DELETE` | `/api/post/{POST_ID}/{COMMENT_ID}` | Удаление комментария            |
| `GET`    | `/api/post/{POST_ID}/upvote`       | Лайк поста                      |
| `GET`    | `/api/post/{POST_ID}/downvote`     | Дизлайк поста                   |
| `GET`    | `/api/post/{POST_ID}/unvote`       | Отмена голосования              |
| `DELETE` | `/api/post/{POST_ID}`              | Удаление поста                  |
| `GET`    | `/api/user/{USER_LOGIN}`           | Посты конкретного пользователя  |
