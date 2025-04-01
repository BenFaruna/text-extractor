build:
	@ go build -C cmd/extractor && cp cmd/extractor/extractor .
clean:
	@ rm cmd/extractor/extractor extractor *.txt
test-pdf:
	@ go test ./pkg/formats/pdf -v
