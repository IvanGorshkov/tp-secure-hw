# tp-secure-hw

## Настройка  (MACOS)
Дальше Настройки > Сеть > Дополнительно > Прокси > http 127.0.0.1:8080 > OK > Применить


## Запуск

```
 go run main.go   
```

или в Docker

```
docker build -t server -f Dockerfile . 
docker run -p 8080:8080  server   
```