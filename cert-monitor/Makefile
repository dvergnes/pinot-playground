
# Build info
BUILD_INFO_IMPORT_PATH=github.com/dvergnes/pinot-playground/cert-monitor/internal/version
VERSION=$(shell git describe --always --match "v[0-9]*" HEAD)
BUILD_INFO=-ldflags "-X $(BUILD_INFO_IMPORT_PATH).Version=$(VERSION)"
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

PKGS=$(shell go list ./...)

.PHONY: gotidy
gotidy:
	rm -fr go.sum
	go mod tidy -go=1.17

.PHONY: gomoddownload
gomoddownload:
	go mod download

.PHONY: install-tools
install-tools:
	brew install mockery

.PHONY: generate-mock
generate-mock: install-tools
	mockery --case underscore --dir alert --name Notifier
	mockery --case underscore --dir monitor --name CertificateInfoGatherer
	mockery --case underscore --dir monitor --name Clock
	mockery --case underscore --name Interface --srcpkg github.com/jetstack/cert-manager/pkg/client/clientset/versioned
	mockery --case underscore --name CertmanagerV1Interface --srcpkg github.com/jetstack/cert-manager/pkg/client/clientset/versioned/typed/certmanager/v1
	mockery --case underscore --name CertificateInterface --srcpkg github.com/jetstack/cert-manager/pkg/client/clientset/versioned/typed/certmanager/v1

.PHONY: cert-monitor
cert-monitor:
	GO111MODULE=on CGO_ENABLED=0 go build -trimpath -o ./bin/cert-monitor_$(GOOS)_$(GOARCH)$(EXTENSION) \
		$(BUILD_INFO) ./cmd

.PHONY: cert-monitor-linux_amd64
cert-monitor-linux_amd64:
	GOOS=linux   GOARCH=amd64 $(MAKE) cert-monitor

.PHONY: cert-monitor-darwin_amd64
cert-monitor-darwin_amd64:
	GOOS=darwin   GOARCH=amd64 $(MAKE) cert-monitor

.PHONY: docker
docker: cert-monitor-linux_amd64
	docker build -t cert-monitor:latest .