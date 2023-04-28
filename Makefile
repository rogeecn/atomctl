.PHONY: install
install:
	rm -f $(GOPATH)/bin/atomctl
	go install .