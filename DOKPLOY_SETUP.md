# Инструкция по настройке Dokploy

## Архитектура приложения

Приложение состоит из трех сервисов:
- **Client** (React + Vite) - веб-интерфейс
- **Server** (Go) - API сервер
- **Database** (PostgreSQL) - база данных

## Схема доменов

- Клиент: `https://yuki.stepan-vinokurov-moscow.ru`
- API: `https://yuki.stepan-vinokurov-moscow.ru/api`

## Шаг 1: Подготовка проекта

1. Убедитесь, что в корне репозитория находится файл `test-deploy.yml`
2. Убедитесь, что все изменения закоммичены и запушены в GitHub

## Шаг 2: Создание проекта в Dokploy

1. Войдите в панель управления Dokploy
2. Создайте новый проект (Project)
3. Выберите тип развертывания: **Docker Compose**
4. Укажите:
   - **Repository**: `https://github.com/S-VIN/yuki_buy_log.git`
   - **Branch**: `main` (или ваша рабочая ветка)
   - **Compose File**: `test-deploy.yml`

## Шаг 3: Настройка доменов и проксирования

### 3.1 Настройка клиентского приложения

1. В Dokploy найдите сервис `client` в вашем проекте
2. Перейдите в раздел **Domains**
3. Добавьте домен:
   - **Domain**: `yuki.stepan-vinokurov-moscow.ru`
   - **Port**: `5173` (внешний порт из test-deploy.yml)
   - **Path**: `/`
   - **SSL**: Включить (Let's Encrypt)

### 3.2 Настройка API сервера

1. В Dokploy найдите сервис `server` в вашем проекте
2. Перейдите в раздел **Domains**
3. Добавьте домен с путем:
   - **Domain**: `yuki.stepan-vinokurov-moscow.ru`
   - **Port**: `8080` (внешний порт из test-deploy.yml)
   - **Path**: `/api`
   - **SSL**: Включить (Let's Encrypt)

### 3.3 Настройка перезаписи путей (Path Rewrite)

В настройках домена для сервера включите **Strip Path Prefix**, чтобы запросы к `/api/*` перенаправлялись на сервер как `/*`.

Пример:
- Запрос: `https://yuki.stepan-vinokurov-moscow.ru/api/products`
- Перенаправляется на сервер как: `http://server:8080/products`

## Шаг 4: Переменные окружения (если требуется изменить)

Все переменные окружения уже прописаны в `test-deploy.yml`, но вы можете их переопределить в интерфейсе Dokploy:

### Для сервиса `server`:
```env
DATABASE_URL=postgres://user:pass@db:5432/yukibuylog?sslmode=disable
SERVER_PORT=8080
CORS_ORIGIN=https://yuki.stepan-vinokurov-moscow.ru
```

### Для сервиса `client`:
Build Args (устанавливаются при сборке):
```env
VITE_API_URL=https://yuki.stepan-vinokurov-moscow.ru/api
```

### Для сервиса `db`:
```env
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=yukibuylog
```

**ВАЖНО**: Для production рекомендуется изменить пароль БД и использовать более надежный секретный ключ!

## Шаг 5: Деплой приложения

1. В интерфейсе Dokploy нажмите **Deploy**
2. Dokploy автоматически:
   - Склонирует репозиторий
   - Соберет Docker образы для всех сервисов
   - Запустит контейнеры
   - Настроит сеть между сервисами
   - Настроит SSL сертификаты

3. Мониторьте процесс деплоя в логах Dokploy

## Шаг 6: Проверка работоспособности

После успешного деплоя:

1. Откройте `https://yuki.stepan-vinokurov-moscow.ru` - должен загрузиться клиент
2. Проверьте API: `https://yuki.stepan-vinokurov-moscow.ru/api/products` - должен вернуть данные или запросить авторизацию

## Автоматический деплой при изменениях

Dokploy поддерживает автоматический деплой через webhooks:

1. В настройках проекта включите **Auto Deploy**
2. Скопируйте webhook URL
3. В настройках GitHub репозитория добавьте webhook:
   - **Payload URL**: скопированный URL из Dokploy
   - **Content type**: `application/json`
   - **Events**: Push events

Теперь каждый push в репозиторий будет автоматически запускать деплой!

## Полезные команды для отладки

Если нужно проверить логи или состояние контейнеров:

```bash
# Просмотр логов сервера
docker logs yuki_server

# Просмотр логов клиента
docker logs yuki_client

# Просмотр логов БД
docker logs yuki_db

# Проверка сети
docker network inspect yuki_buy_log_yuki_network
```

## Troubleshooting

### Проблема: API недоступно с клиента

- Проверьте CORS_ORIGIN в настройках сервера
- Убедитесь, что path rewrite настроен корректно
- Проверьте логи сервера на предмет CORS ошибок

### Проблема: Сервер не может подключиться к БД

- Убедитесь, что все сервисы в одной сети (`yuki_network`)
- Проверьте healthcheck БД
- Проверьте DATABASE_URL в настройках сервера

### Проблема: Клиент показывает неправильный API endpoint

- Пересоберите клиент с правильным `VITE_API_URL`
- Убедитесь, что build arg передается при сборке образа

## Безопасность

Для production рекомендуется:

1. Изменить пароль БД в переменных окружения
2. Использовать Docker secrets для чувствительных данных
3. Настроить rate limiting на уровне прокси
4. Регулярно обновлять зависимости и образы
5. Настроить бэкапы БД
