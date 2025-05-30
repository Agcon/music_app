# 🎵 Music App

Веб-приложение для загрузки, прослушивания и управления музыкальными треками с поддержкой авторизации, аутентификации и хранения файлов в MongoDB через GridFS.

## 🚀 Быстрый старт

### 1. Клонирование репозитория

<details>
<summary>HTTPS</summary>

```bash
git clone https://github.com/Agcon/Music-App.git
```
</details>

<details>
<summary>SSH</summary>

```bash
git clone git@github.com:Agcon/Music-App.git
```
</details>

### 2. Запуск приложения

Перейдите в директорию проекта и выполните:

```bash
cd Music-App
docker-compose up --build
```

Приложение будет доступно по адресу: [http://localhost:8086](http://localhost:8086)

## 🧪 Проверка функциональности

1. Откройте [http://localhost:8086](http://localhost:8086)
2. Зарегистрируйте нового пользователя на странице регистрации.
3. После входа вы попадёте в список треков.
4. В списке можно:
    - 🔊 прослушивать треки,
    - ⬆️ загружать свои треки,
    - 🗑 удалять треки.
5. При первом запуске автоматически загружаются примеры треков из папки `assets`.
6. В верхней части страницы отображается почта текущего пользователя и доступна кнопка выхода.

---

Проект использует PostgreSQL, MongoDB, Redis и написан на Go с использованием Gin и GORM. Все зависимости разворачиваются автоматически через Docker Compose.

