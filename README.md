## Превьювер изображений

Для локальной разработки кэшируемые изображения складываются в папку `cache`. Размер кэша измеряется в кол-ве закэшированных изображений.

### DONE

* HTTP-сервер, проксирующий запросы к удаленному серверу
* Докер и Makefile
* Нарезка изображений
* Кэширование нарезанных изображений на диске
* Ограничение кэша кол-вом изображений
* Тесты кэша
* Интеграционные тесты

### Запуск в docker

```
make build-docker
make run
```

### Запуск тестов

```
make test
```

### Запуск линтера

```
make lint
```
