перенести миграции в папку migrations 
при раскатке сервиса (docker run), нужно каждому сервису чекать (мэйн версия для всех)

ручку для ордера (post для кастомера) get для курьера и для ресторана (в pkg)

сначала только авторизацию наверн по образу и подобию для других

---

kibana прометус кафка логброкер rest api кубер

линукс в очко

были тайминги 100 милискину стали 20

mau с 3.000.000 до 6.000.000
dau 

прогуглить продукт метрики

залупа

улучшил тайминги

была ручка что то то делала за 150мс а я нахуяервитл за 70 но надо за 20

 Покрывал код unit- и интеграционными тестами с testify и testcontainers, участвовал в
командных код-ревью -- убрать

вот такое flow тестирование 

end-to-end

 Настроил сбор метрик в Prometheus, построил дашборды в Grafana, добавил распределённый
трейсинг через Jaeger. ОСТАВИТЬ


НЕ Jaeger а эластик


Дима на вб был в банкинге Надо было раработать систему стимуляции юзеров 

тегнуть Диму голосовое 




# Запуск с данными покупателей (100 итераций)
npx newman run YAFDS.postman_collection.json -e YAFDS.postman_environment.json --iteration-data migrations/testdata/customer/postman_customers_true_100.csv

# Запуск с данными курьеров
npx newman run YAFDS.postman_collection.json -e YAFDS.postman_environment.json --iteration-data migrations/testdata/courier/postman_couriers_true_100.csv

# Запуск с заказами
npx newman run YAFDS.postman_collection.json -e YAFDS.postman_environment.json --iteration-data migrations/testdata/orders/postman_orders_true_100.csv