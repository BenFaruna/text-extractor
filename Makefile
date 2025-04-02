build:
	@ go build -C cmd/extractor && cp cmd/extractor/extractor .
clean:
	@ rm cmd/extractor/extractor extractor *.txt
test-pdf:
	@ go test ./pkg/formats/pdf -v
test-text:
	@ go test ./pkg/formats/text -v
test-doc:
	@ go test ./pkg/formats/doc -v