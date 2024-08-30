# Catbox clone

Accepts files upload on /, protected by HTTP basic auth. Works through curl too. With admin/password as username/password:

```sh
$ curl -u admin:password -X POST -F file=@coolfile.zip http://localhost:8080
```

Needs a `config.json` with the username and password like so:

```json
{
		"username": "admin",
		"password": "password"
}
```

