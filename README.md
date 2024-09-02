# Описание проекта
- Проект сделан, следуя дизайну REST API.
- Проект запускается из Docker.
- База данных Postgresql. Подключение к базе данных с помощью <a href="https://github.com/jmoiron/sqlx">sqlx</a>
- Регистрация и аутентификация. Работа с JWT. Middleware.
- Тестирование с помощью <a href="https://github.com/stretchr/testify">testify</a> и <a href="https://github.com/golang/mock">gomock</a>
- Добавлены метрики Prometheus с отображением через Grafana


# Build
Run app in docker

    make run


Stop app in docker

    make stop
    
Run tests

    make test

# Info

Swagger available on http://localhost:8080/swagger

pgAdmin available on http://localhost:5050

```
login: admin@admin.com
password: admin

postgres_host: db
postgres_database: foxgres
postgres_password: admin
```

Grafana available on http://localhost:3000/

Можно импортировать дашборд [Percentage posts.json](grafana%2FPercentage%20posts.json)

## TODO
Можно сделать, но не успел:
- Graceful Shutdown
- Тестирование `storage` и `service`