include .envrc

.PHONY: build
build:
	docker build -t ${AV_SCANNER_IMAGE_TAG} .

.PHONY: run
run: build
	docker run --rm -it \
      -v ${AV_SCANNER_CLAMAV_DEFINITION_DIR}:/var/lib/clamav:rw \
      -v ${AV_SCANNER_SCAN_DIR}:/temp/scan:rw \
      -p ${AV_SCANNER_PORT}:8080 \
      ${AV_SCANNER_IMAGE_TAG}