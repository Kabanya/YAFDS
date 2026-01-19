##### courier
0x5647445a3e564b189fc35a75079fc924 aad9D6C7e9R7
0xc56289d4aabb42149853763cfa77262d cdS8ZrCsdp6Q

##### customer
0x742d35Cc6634C0532925a3b844Bc454e4438f44e StrongPass001
0x1f9840a85d5af5bf1d1762f925bdaddc4201f984 MySecret003

##### restaurant
0xbd2f054c1faa44b2a561c47dd1a0a367 JO88M0Rtf0BG   ee4ef0d3-7ee3-4cdc-9adb-cd36a8971c8a
0x851db72d3ba64acab39cd825ed82dea6 MqFhwJG8N6NO



в нашем корпоративном? мире в pkg тоже самое что и в основном (структура)

запросить доку репозитория


1. res menu item 0 сначала сделать (в сервисе с рестораном)
2. расширить создания ордера чтобы он дергал рид меню айтем
3. ручка в текущий ордер добавит ресторан меню айтем
4. эта ручка чекает айтем доступен (quant > 0)
5. ttl должен быть у вещей ресторана. если условно > 10 минут то <новое>

на стороне заказчика бросить ивент (создан новый заказ) в кафку. Кто то на стороне рестика читает и не/принимает
kafka


accept order проходит по quants меню айтемс и если > n , то true(order_state = kitchen_accepted) else false (...denied)

или ручка или usecase на стороне customers - PAY : запросить у wallet_address есть ли деньги и если да, то списать. Перевести в статус customer_paid else customer_cancelled

нужен в pkg отдельный usecase статус изменения order'а. Он чекает порядок stat'ов. По древу состояний двигаемся только вниз (order states.webp)

---

end-to-end (e2e) тесты как душе угодно (но лучше всего на том же языке). Делаем на уровне repository папку test. test_order_flow.go передаем нужные ручки и тд и тп. Чекаем верхноуровнево.
для них полезно делать в postman коллекции

unit-тесты каждый слой тестируют. Тестируем маленький кусочек кода и все внешние зависимости кода mock'аем (go-mock тулза). mock = заглушка, условно всегда на запрос true возрващаем (https://pkg.go.dev/github.com/golang/mock/gomock#example-Call.DoAndReturn-Latency)
Нужен код который мы дернем из customer и который под копотом дернет http ресторана. Customer/clients (clients = сервисы в которые мы ходим).

---

!!! всё что пишется покрываем тестами.

end-to-end ручки ордера, что она что то создает и дергает
List Orders чекает orders если пусто то create order.

unit|mock test на логику usecase добавления order'а (строчка 20)

---

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