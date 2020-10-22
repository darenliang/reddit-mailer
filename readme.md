# Reddit Mailer

Reddit Mailer is a small cross-platform program that sits neatly in your system-tray and checks your Reddit mail!

![](https://raw.githubusercontent.com/George-lewis/reddit-mailer/master/screenshots/nomail-windows.png)

![](https://raw.githubusercontent.com/George-lewis/reddit-mailer/master/screenshots/mail-windows.png)

![](https://raw.githubusercontent.com/George-lewis/reddit-mailer/master/screenshots/nomail-ubuntu.png)

![](https://raw.githubusercontent.com/George-lewis/reddit-mailer/master/screenshots/mail-ubuntu.png)

## Acquiring It

If you're on Windows or Linux go ahead and grab an artifact from the Github Actions tab then go to `Using It`.

Or if want to build, run `go build -o reddit-mailer main.go` on linux and run the batch file if you're on Windows.

If you're on MacOS the situation is more complicated, please read: https://github.com/getlantern/systray#macos

## Using It

You'll need to supply an agent file for the program to work, the program expects this file to be called `agent.txt` and placed in the same folder as the executable. [This wiki page](https://turnage.gitbooks.io/graw/content/chapter1.html) explains how to create an agent file.

After that place the executable, icons, and your agent.txt all in the same place and run the program!

## Configuring It

Reddit Mailer comes with a configuration file `config.json`, this file must be present in the same directory as the executable

> `limit; int`

This value represents the number of unread mail retrieved by the application. If you ever have >= limit unread messages, the application won't count them and will display "limit+". e.g. if you have >=30 unread messages (the default limit) the systray will read: "30+ messages"

> `interval; int`

Check the inbox every `interval` seconds.

Note: The application also checks the inbox when it receives an event from the API

> `notifications; bool`

This enables or disables notifications, takes precedence over toggles for specific types

> `comment_replies, messages, post_replies, mentions; bool`

These are specific toggles for notifications regarding a specific type of unread message