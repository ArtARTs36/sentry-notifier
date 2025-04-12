# sentry-notifier

**sentry-notifier** - приложение для отправки уведомлений по веб-хукам от Sentry

Поддерживаемые мессенджеры:
- Telegram
- Mattermost

Приложение поднимает 2 порта:
- 8080 (main server) - для принятия веб-хуков от Sentry
- 8081 (control server) - ручки для отправки тестовых уведомлений, метрик и Health Check

