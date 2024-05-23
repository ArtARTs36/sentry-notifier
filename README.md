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
  client_secret: '$SENTRY_NOTIFIER_SENTRY_CLIENT_SECRET'

channels:
  my_team:
    telegram:
      - chat_id: '$SENTRY_NOTIFIER_TELEGRAM_CHAT_ID'
        thread_id: '$SENTRY_NOTIFIER_TELEGRAM_CHAT_THREAD_ID'
        bot_token: '$SENTRY_NOTIFIER_TELEGRAM_BOT_TOKEN'

notify:
  strategy: async
  on:
    event_alert:
      - message: |
          🚨 Error on {{ hook.Extracted.OrganizationName }}/{{ hook.Extracted.ProjectName }}

          at {{ hook.Event.Datetime.Human() }}

          ```
          {{ hook.Event.Title }}```

          {{ hook.Event.IssueURL }}
        to: my_team

    issue:
      - message: |
          🚨 Error on {{ hook.Issue.Project.Name }}

          at {{ hook.Issue.LastSeen.Human() }}

          ```
          {{ hook.Issue.Title }}```

          https://<your-org-name>.sentry.io/issues/{{ hook.Issue.ID }}/
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
    environment:
      - SENTRY_NOTIFIER_SENTRY_CLIENT_SECRET=<your sentry client secret>
      - SENTRY_NOTIFIER_TELEGRAM_CHAT_ID=<your telegram chat id>
      - SENTRY_NOTIFIER_TELEGRAM_CHAT_THREAD_ID=<your telegram chat thread id>
      - SENTRY_NOTIFIER_TELEGRAM_BOT_TOKEN=<your telegram bot token>
```