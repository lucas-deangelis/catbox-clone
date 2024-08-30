copy-to-vps: build
	rsync -avz catbox-clone.service catbox-clone manticore:catbox-clone/

build:
	go build

