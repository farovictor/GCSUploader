init-gcs-emulator:
	docker run -d \
		-e PORT=9023 \
		-p 9023:9023 \
		--name gcp-storage-emulator \
		--rm \
		oittaa/gcp-storage-emulator

build:
	go build -o ./bin/gcsloader ./src
	chmod a+x ./bin/gcsloader

run-create-bucket: build
	./bin/gcsloader create-bucket --bucket-name test --project-id local

run-delete-bucket: build
	./bin/gcsloader delete-bucket --bucket-name test

run-check-bucket-exists: build
	./bin/gcsloader exists-bucket --bucket-name test

run-check-bucket-attrs: build
	./bin/gcsloader attrs-bucket --bucket-name test

run-load: build
	./bin/gcsloader load --bucket-name test --blob-path "some/weird/path/blob.txt" --source-file "Makefile"

run-load-prefix: build
	./bin/gcsloader load-prefix --bucket-name test --search-path "./data" --blob-prefix-path "g/prefix" --blob-prefix-name "" --num-concurrent-files 50


.PHONY: init-gcs-emulator