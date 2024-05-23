# sentry-notifier

**sentry-notifier** - app for sending notifications on Sentry webhooks

Supported messengers:
- Telegram

## Configuration

Configuration structure declared in [sentry-notifier-json-schema.json](./sentry-notifier-json-schema.json)

```yaml
http:
  address: :8088

security:
  client_secret: '$SENTRY_CLIENT_SECRET'

channels:
  my_team:
    telegram:
      - chat_id: '$TELEGRAM_CHAT_ID'
        thread_id: '$TELEGRAM_CHAT_THREAD_ID'
        bot_token: '$TELEGRAM_BOT_TOKEN'

notify:
  strategy: async
  on:
    event_alert:
      - message: |
          🚨 Error on {{ hook.Extracted.OrganizationName }}/{{ hook.Extracted.ProjectName }}
          
          at {{ hook.Event.Datetime.Human() }}

          ```
          {{ hook.Event.Title }}```
          
          {{ hook.Event.WebURL }}
        to: my_team
```

## Run with docker-compose

```yaml
services:
  sentry-notifier:
    image: artarts36/sentry-notifier:0.1.0
    ports:
      - "80:8088"
    volumes:
      - ./sentry-notifier.yaml:/app/sentry-notifier.yaml
```