SMTP2HTTP (email-to-web)
========================
smtp2http is a simple smtp server that resends the incoming email to the configured web endpoint (webhook) as a basic http post request.

Dev 
===
- `go mod vendor`
- `go build`

Dev with Docker
==============
Locally :
- `go mod vendor`
- `docker build -f Dockerfile.dev -t mail2gotify-proxy-dev .`
- `docker run -p 25:25 mail2gotify-proxy-dev --timeout.read=50 --timeout.write=50 --gotify https://gotify.example.com`

Or build it as it comes from the repo :
- `docker build -t mail2gotify-proxy .`
- `docker run -p 25:25 mail2gotify-proxy --timeout.read=50 --timeout.write=50 --gotify https://gotify.example.com`

use application name as user and apptocken as password when sending Mails
Append `-<priority>` to user to set the priority of all messages

The `timeout` options are of course optional but make it easier to test in local with `telnet localhost 25`
Here is a telnet example payload : 
```
HELO zeus
# smtp answer

AUTH PLAIN <<base64 encoded user/password>
# smtp answer

MAIL FROM:<email@from.com>
# smtp answer

RCPT TO:<youremail@example.com>
# smtp answer

DATA
your mail content
.

```

Docker (production)
=====
**Docker images arn't available online for now**
**See "Dev with Docker" above**
- `docker run -p 25:25 smtp2http --gotify https://gotify.example.com --listen :2525`

Native usage
=====
`smtp2http --listen=:25 --webhook=http://localhost:8080/api/smtp-hook`
`smtp2http --help`

Contribution
============
Original repo from @alash3al
Thanks to @aranajuan


