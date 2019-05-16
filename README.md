# go-telenote
Send simple telegram notifications

Usage:
```
botToken := "your-bot-token"
chatId := -12345678
text := "Hello world"
notifier = telenote.NewNotifier(botToken)

err := notifier.Notify(chatId, text)

// or using a custom http client
err := notifier.WithClient(&http.Client{}).Notify(chatId, text)
```
