# Makefile for this project

PACKAGE_ROOT=github.com/markmandel/recommendation-neo4j
NEO4J_DATA=./neo4j-data
BACKUP_FILE=data-backup.7z

.PHONY: rescue test deviation

#start it up
up:
	sudo fig up

#install all the binaries
install:
	go install github.com/golang/lint/golint
	go install $(PACKAGE_ROOT)/seed

rescue:
	go build -o ./bin/rescue $(PACKAGE_ROOT)/rescue
	./bin/rescue

deviation:
	go build -o ./bin/deviation $(PACKAGE_ROOT)/deviation
	./bin/deviation

#check the code for correctness
check-code:
	go vet $(PACKAGE_ROOT)/...; true
	golint $(PACKAGE_ROOT)/...; true

#gets the doc server up and running
doc:
	killall godoc; godoc -http=":7080" &

#get all the deps
deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/PuerkitoBio/goquery
	go get -u github.com/manki/flickgo
	go get -u github.com/jmcvetta/neoism
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/sessions
	go get -u github.com/cenkalti/backoff

neo4j-clean:
	sudo fig stop
	sudo rm -rf $(NEO4J_DATA)/graph.db

neo4j-backup:
	sudo fig stop
	rm $(NEO4J_DATA)/$(BACKUP_FILE); true
	7z a $(NEO4J_DATA)/$(BACKUP_FILE) ./neo4j-data/graph.db

neo4j-restore:
	$(MAKE) neo4j-clean
	cd $(NEO4J_DATA) && 7z x $(BACKUP_FILE)