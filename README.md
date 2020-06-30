# Arveto Auth

(public package) [![GoDoc](https://godoc.org/github.com/Arveto/arvetoAuth/pkg/public?status.svg)](https://godoc.org/github.com/Arveto/arvetoAuth/pkg/public)

## Start

Start this programme with the name of the working directory.

## Config

```ini
; URL of the server
url = "http://localhost:8000/"
; Listen address
listen = ":8000"

; Enable develpment mode
; ! NO USE IT IN PRODUCTION !
dev = true

[github]
; The GitHub Oauth tocken
client = ""
secret = ""

[google]
; The Google Oauth tocken
client = ""
secret = ""

[mail]
; Mail configuration
login = "auth@example.com"
password = "xxx"
host = "smtp.example.com"
```

## JWT

You can get the RSA public key from `/publickey` URI.

The JWT claims:

-   `id` (string): an string id.
-   `pseudo` (string)
-   `email` (string): email address.
-   `avatar` (string): The URL of teh avatar or a default avatar.
-   `level` (string): The level of the user. Values are `Ban` &lt; `Candidate` &lt; `Visitor` &lt; `Std` &lt; `Admin`.

To auth a client, redirect it to `http(s)://server/auth?app=appID&r=base64`.
The params r is an optionnal redirect URL encoding in URL base64. The client
come back on the URI: `/login?jwt=...?r=...`.
