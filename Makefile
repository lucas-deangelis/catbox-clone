copy-to-vps: build
	rsync -avz catbox-clone.service catbox-clone config.json manticore:catbox-clone/

build:
	go build

