SHELL = /bin/zsh

include .env

export BUCKET_NAME
export JSON_CREDS
export PROJECT_ID
export STORAGE_EMULATOR_HOST

init-gcs-emulator:
	docker run -d \
		-e PORT=9023 \
		-p 9023:9023 \
		--name gcp-storage-emulator \
		--rm \
		oittaa/gcp-storage-emulator

build:
	@go build --race -o ./bin/gcsloader ./src
	@chmod a+x ./bin/gcsloader

run-create-bucket: build
	@./bin/gcsloader -a "emulator" create-bucket --bucket-name "$(BUCKET_NAME)" --project-id "$(PROJECT_ID)"

run-delete-bucket: build
	@./bin/gcsloader -a "emulator" delete-bucket --bucket-name "$(BUCKET_NAME)"

run-check-bucket-exists: build
	@./bin/gcsloader -a "emulator" exists-bucket --bucket-name "$(BUCKET_NAME)"

run-check-bucket-attrs: build
	@./bin/gcsloader -a "emulator" attrs-bucket --bucket-name "$(BUCKET_NAME)"

run-load: build
	@./bin/gcsloader -a "emulator" load --bucket-name "$(BUCKET_NAME)"--blob-path "some/weird/path/blob.txt" --source-file "Makefile"

run-load-in-batches: build
	@./bin/gcsloader -a "emulator" load-batch --bucket-name "$(BUCKET_NAME)" --search-path "./data" --blob-prefix-path "some/prefix/now_100_workers" --blob-prefix-name "" --num-concurrency 100

run-load-gcp-in-batches: build
	@./bin/gcsloader -a "json" -c "$(JSON_CREDS)" load-batch --bucket-name "$(BUCKET_NAME)" --search-path "./data" --blob-prefix-path "some/prefix/now_100_workers" --blob-prefix-name "" --num-concurrency 100

run-print-creds:
	@echo "my bucket: $(BUCKET_NAME)"
	@echo "there is a json content: $(JSON_CREDS)"

.PHONY: init-gcs-emulator, run-create-bucket, run-delete-bucket, run-check-bucket-exists, run-check-bucket-attrs, run-load, run-load-in-batches, run-load-gcp-in-batches, run-print-creds
