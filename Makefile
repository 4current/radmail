target = radmail
gopath = ${HOME}/go
gopkg = github.com/4current/radmail

build:
	go build $(gopkg)
run:
	go run $(gopkg)
install:
	go install $(gopkg)
uninstall:
	rm $(gopath)/bin/$(target)
