# sentry-notifier

**sentry-notifier** - app for sending notifications on Sentry webhooks

Supported messengers:
- Telegram
- Mattermost

## Configuration

Configuration structure declared in [sentry-notifier-json-schema.json](./sentry-notifier-json-schema.json)

```yaml
http:
  address: :8080
  
control:
  address: :8081

security:
  client_secret: '$SENTRY_NOTIFIER_SENTRY_CLIENT_SECRET'

channels:
  my_team:
    telegram:
      - chat_id: '$SENTRY_NOTIFIER_TELEGRAM_CHAT_ID'
        thread_id: '$SENTRY_NOTIFIER_TELEGRAM_CHAT_THREAD_ID'
        bot_token: '$SENTRY_NOTIFIER_TELEGRAM_BOT_TOKEN'
        
    mattermost:
      - server: 'http//localhost:8065'
        token: '$SENTRY_NOTIFIER_MATTERMOST_TOKEN'
        channel:
          name: 'alerts'
          team_name: 'My Team'

notify:
  strategy: async
  on:
    event_alert:
      - message: |
          🚨 Error on {{ hook.Extracted.OrganizationName }}/{{ hook.Extracted.ProjectSlug }}

          at {{ hook.Event.Datetime.Human() }}

          ```
          {{ hook.Event.Title }}```

          {{ hook.Event.WebURL }}
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

## Templating

[Twig syntax](https://twig.symfony.com) is used to compile templates. The Twig port is a [Stick](https://github.com/tyler-sommer/stick).

**event_alert**

| Variable                    | Type     | Description                                              |
|-----------------------------|----------|----------------------------------------------------------|
| hook.Event.IssueURL         | string   | The API URL for the associated issue                     |
| hook.Event.IssueID          | string   | The id of the issue                                      |
| hook.Event.Platform         | string   |                                                          |
| hook.Event.Title            | string   | The label of the rule that was triggered                 |
| hook.Event.Type             | string   |                                                          |
| hook.Event.Project          | integer  |                                                          |
| hook.Event.URL              | string   |                                                          |
| hook.Event.Datetime         | time     |                                                          |
| hook.Event.Datetime.Human() | string   | Format time to `Y-m-d H:i:s`                             |
| hook.Event.URL              | string   | The API URL for the event                                |
| hook.Event.Fingerprint      | string[] |                                                          |
| hook.Event.Request.Method   | string   |                                                          |
| hook.Event.Request.URL      | string   |                                                          |
| hook.Extracted.ProjectName  | string   | Name of your project, extracted from hook.Event.URL      |
| hook.Extracted.Organization | string   | Name of your organization, extracted from hook.Event.URL |

**issue**

| Variable                     | Type   | Description                                                           |
|------------------------------|--------|-----------------------------------------------------------------------|
| hook.Issue.Count             | string |                                                                       |
| hook.Issue.ID                | string |                                                                       |
| hook.Issue.Action            | string | can be `created`, `resolved`, `assigned`, `archived`, or `unresolved` |
| hook.Issue.Level             | string |                                                                       |
| hook.Issue.ShortID           | string |                                                                       |
| hook.Issue.Status            | string |                                                                       |
| hook.Issue.Type              | string |                                                                       |
| hook.Issue.Title             | string |                                                                       |
| hook.Issue.FirstSeen         | time   |                                                                       |
| hook.Issue.LastSeen          | time   |                                                                       |
| hook.Issue.FirstSeen.Human() | string | Format time to `Y-m-d H:i:s`                                          |
| hook.Issue.LastSeen.Human()  | string | Format time to `Y-m-d H:i:s`                                          |
| hook.Issue.Project.ID        | string |                                                                       |
| hook.Issue.Project.Name      | string |                                                                       |
| hook.Issue.Project.Platform  | string |                                                                       |
| hook.Issue.Project.Slug      | string |                                                                       |
