# tp-secure-hw

 ##Настройка  (MACOS)
Запуск gen_ca
```
./gen_ca.sh
```
Файл ca.cert поместить в связку ключей. В параметрах сертификата указать "Доверять всегда".

Дальше Настройки > Сеть > Дополнительно > Прокси > http/https 127.0.0.1:8080 > OK > Применить


##Запуск

```
 go run main.go   
```

или в Docker

```
docker build -t server -f Dockerfile . 
docker run -p 8080:8080  server   
```