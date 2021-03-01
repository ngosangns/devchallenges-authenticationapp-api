build:
	set GOOS=linux
	go build -o functions/hello hello.go
dev:
	netlify dev
deploy:
	netlify deploy --dir=site --functions=functions --prod