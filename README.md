# hexkey

*(Feel free to propose a better name)*

Hexkey generates a hexagon from a 256 bit value encoded using base64 (i.e. public keys and hashes).

```
git clone ssb://%1jSkm2ziiZ9FbO/kyhRyd3gtn9UtbJQLzYf13HgRO4E=.sha256 hexkey
cd hexkey
go build
./hexkey $(sbot whoami | grep -o @.*\.ed25519sbot whoami | grep -o @.*\.ed25519)
```
