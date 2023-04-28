# wbL0
test task L0 for wildberries

# Схема данных для PostgreSQL
На типы данных можно посмотреть в миграции /internal/db/migrations/001_add_users.up.sql
Касаемо связи между таблицами:
1) orders - orders_uid primary key сделал ему тип uuid, чтобы продемонстрировать умение работать с uuid
2) deliveries - создано bigint поле - primary key с названием delivery_id, внешний ключ к таблице orders - orders_uid (unique т.к. для одного заказа одна доставка) 
3) payments - создано bigint поле - primary key с названием payment_id, внешний ключ к таблице orders - transaction (unique т.к. для одного заказа одна оплата)
4) item - создано bigint поле - primary key с названием item_id, внешний ключ к таблице orders - orders_uid (для одного заказа может быть много item'ов)

Касаемо остальных данных: у каждого item имеется track_number, которые совпадает с track_number у order, но я не стал уделять этому должно внимания и делать проверки на совпадения, потому что primary key у order уже есть, а требования по реализации бизнес-логики, связанной с данными не было.

# Скрипт для записи данных в канал
Репозиторий: https://github.com/AlexeyBazhin/wbL0Script
Записывает за одно выполнение по 5 рандомно сгенирированых на основе шаблона замаршаленных в виде слайса байт JSON в кластер nats-streaming.

# В целом о проекте
Реализована clean архитектура. На уровне api имеется 2 параллельно работающие сущности (http server и stanListener). Они вызываются с помощью errGroup в main, а также разделяют общий service (и следовательно repository).
В качестве кэша использовался Redis
В проекте поднимаются docker контейнеры с: postgres, stan, redis и самим приложением 

# "Субтитры" к видео
