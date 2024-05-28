sudo systemctl stop showmaster.service

git pull

go build ./src/main.go

sudo systemctl start showmaster.service