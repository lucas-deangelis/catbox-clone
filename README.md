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

## How do I install this wonderful program on my VPS that runs Ubuntu and just happens to be named "manticore" in my SSH config, just like in your Makefile

First of all make sure you have a `config.json` in the current folder, see example above.

Assuming it's x86 (if it's ARM you'll have to change the `build` target of the makefile to build for ARM):

```sh
$ make copy-to-vps
$ ssh manticore
$ cd catbox-clone
$ sudo cp catbox-clone.service /etc/systemd/system
$ sudo systemctl enable catbox-clone.service
$ sudo systemctl start catbox-clone.service
```

## Awesome. Now how do I serve that on files.example.com? I just happen to use Caddy

Add the following to your Caddyfile (`sudoedit /etc/caddy/Caddyfile`):

```caddyfile
files.example.com {
	reverse_proxy http://localhost:8080
}
```

## Why are all those instructions so specific?

If you're reading it, it's for you.

