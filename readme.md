# Reddit Mailer

Reddit Mailer is a small cross-platform program that sits neatly in your system-tray and checks your Reddit mail!

![](https://raw.githubusercontent.com/George-lewis/reddit-mailer/master/screenshots/nomail.png)

![](https://raw.githubusercontent.com/George-lewis/reddit-mailer/master/screenshots/mail.png)

# Acquiring It

If you're on Windows or Linux go ahead and grab an artifact from the Github Actions tab then go to `Using It`.

Or if want to build, run `go build .` on linux and run the batch file if you're on Windows.

If you're on MacOS the situation is more complicated, please read: https://github.com/getlantern/systray#macos

# Using It

You'll need to supply an agent file for the program to work, the program expects this file to be called `agent.txt` and placed in the same folder as the executable. [This wiki page](https://turnage.gitbooks.io/graw/content/chapter1.html) explains how to create an agent file.

After that place the executable, icons, and your agent.txt all in the same place and run the program!