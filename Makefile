.ONESHELL:

# ANSI color definition
BgLightRed := \033[101m
BgLightGreen := \033[102m
BgLightBlue := \033[104m
FgBlack := \033[30m
FgRed := \033[31m
FgGreen := \033[32m
FgLightBlue := \033[94m
ColorClose := \033[0m

# Scoped variable that can be modified with any
# Others
GO ?= go
GO_INSTALL? = $(GO) install
MIN_SUP_GO_MJR_VER = 1
MIN_SUP_GO_MNR_VER = 16
OK_PREFIX = * [OK..]
FAIL_PREFIX = * [FAIL]
GO_VERSION_VALIDATION_ERR_MSG = Golang version is not supported. Please update to at least v$(MIN_SUP_GO_MJR_VER).$(MIN_SUP_GO_MNR_VER)
# System
PACKAGE ?= $(shell go list)
VERSION ?= $(shell git describe --tags 2> /dev/null || cat pwd/.version 2> /dev/null || echo v0.0.0)
SERVICE_IDENTIFIER ?= $(shell basename `$(GO) list`)
LAST_COMMIT_ID := $(shell git rev-parse HEAD | cut -c 1-12)
LINUX_UNIT_FILE ?= $(SERVICE_IDENTIFIER).service
LINUX_SYSLOG_CONF_FILE ?= syslog-$(SERVICE_IDENTIFIER).conf
LINUX_pId_FILE ?= $(SERVICE_IDENTIFIER).pId
ifeq ($(VERSION), v0.0.0)
VERSION := $(addsuffix -$(shell date "+%Y%m%d%H%M%S")-$(LAST_COMMIT_ID),$(VERSION))
endif
# Golang compile time variable (ldflags)
LD_VAR ?= app.Version app.LastBuildAt
# Golang debug options
DEBUG_OPTIONS ?=

# Scoped constant
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
GO_OS := $(shell go env GOOS)
GO_MAJOR_VERSION := $(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION := $(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
VERSION_GT_REQ := $(shell [ $(GO_MAJOR_VERSION) -gt $(MIN_SUP_GO_MJR_VER) -o \( $(GO_MAJOR_VERSION) -eq $(MIN_SUP_GO_MJR_VER) -a $(GO_MINOR_VERSION) -ge $(MIN_SUP_GO_MNR_VER) \) ] && echo true)
ifeq ($(OS),linux)
DATE := $(shell date --iso=seconds)
else
ifeq ($(OS),darwin)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
endif
endif

# Helper function
define print_error
	"$(BgLightRed)$(FgBlack) ERROR $(ColorClose) $(FgRed)$(1)$(ColorClose)"
endef
define print_success
	"$(BgLightGreen)$(FgBlack) SUCCESS $(ColorClose) $(FgGreen)$(1)$(ColorClose)"
endef
define print_info
	"$(BgLightBlue)$(FgBlack)  $(2)  $(ColorClose) $(FgLightBlue)$(1)$(ColorClose)"
endef

# Doctor variable
# check_pass=0
ifeq ($(OS),$(GO_OS))
OS_OK = 1
endif

## help		: Shows you all available commands
.PHONY: help
help: Makefile
	@echo " Usage: make [target]... [debug options]"
	@echo ""
	@echo " Targets:"
	@sed -n 's/^##//p' $<
	@echo ""
	@echo " Debug Options: (Please start with GODEBUG= and following with golang debug env val, like this)"
	@echo " Please start with GODEBUG= and following with golang debug env val, like this: (i.e. GODEBUG=allocfreetrace=1;schedtrace=1000)"
	@echo " Options Example:"
	@echo " Mem Alloc Trace: allocfreetrace=1"
	@echo " Scheduler Trace: schedtrace=1000"

## doctor		: Checks your environment and display a report to the terminal window.
.PHONY: doctor
doctor:
# Print service information
	@echo $(call print_info,"Service Identifier: $(SERVICE_IDENTIFIER)","INFO")
	@echo $(call print_info,"Package: $(PACKAGE)",\033[8mINFO)
	@echo $(call print_info,"App Version: $(VERSION)",\033[8mINFO)
	@echo $(call print_info,"Current Date: $(DATE)",\033[8mINFO)
	
	@echo ""
	@echo "Diagnosing your system, please wait..."
	@echo ""
	@sleep 2
# System OS vs GOOS
ifneq ($(OS_OK), 1)
	@echo $(call print_error,"OS \(missmatch, GOOS = $(GO_OS)\)")
	exit 1
else
# $(eval check_pass=1)
	@echo $(call print_success,"OS \(your os already same with Go OS env. variable. GOOS = $(GO_OS)\)")
endif
# Current Go Version vs Min. Supported Go Version
ifneq ($(VERSION_GT_REQ), true)
# $(eval check_pass=0)
	@echo $(call print_error,"Golang version \($(GO_VERSION_VALIDATION_ERR_MSG)\)")
	exit 1
else
# $(eval check_pass=1)
	@echo $(call print_success,"Golang version \(is supported \>=$(MIN_SUP_GO_MJR_VER).$(MIN_SUP_GO_MNR_VER)\)")
endif

## genproto	: Compile and generate proto file
.PHONY: genproto
genproto:
	@echo $(shell protoc --go-grpc_out=.public/go --go_out=.public/go --grpc_out=.public/php --php_out=.public/php --plugin=protoc-gen-grpc=.blueprint/grpc_php_plugin .blueprint/insys/*.proto .blueprint/provsys/*.proto)

## genprotojs	: Compile and generate proto file for js defined language
.PHONY: genprotojs
genprotojs:
	@echo $(shell grpc_tools_node_protoc --proto_path=.blueprint --grpc_out=.public/js --js_out=import_style=commonjs,binary:.public/js .blueprint/insys/*.proto .blueprint/provsys/*.proto)

## install	: Install the service into the system.
.PHONY: install
install:
ifeq ($(GO_OS), linux)
install: install_linux install_bin
else
ifeq ($(GO_OS), darwin)
install: install_darwin install_bin
endif
endif

## install_bin	: Build and install app binary in the system (into GOPATH).
install_bin:
	$(shell $(GO_INSTALL))

## update		: Updating the service and app binary (like install, but not copy resources/config files into the system).
.PHONY: update
update:
ifeq ($(GO_OS), linux)
update: update_linux install_bin
else
ifeq ($(GO_OS), darwin)
update: update_darwin install_bin
endif
endif

## smooth_update	: Like update but not restart systemd service (smooth restart handled by overseer package).
.PHONY: smooth_update
smooth_update:
ifeq ($(GO_OS), linux)
smooth_update: smooth_update_linux install_bin
soft_maintenance: soft_maintenance_linux install_bin
hard_maintenance: hard_maintenance_linux install_bin
else
ifeq ($(GO_OS), darwin)
smooth_update: smooth_update_darwin install_bin
endif
endif

# Checker
# checker:
# ifeq ($(check_pass),1)
# 	@echo "Start installing..."
# else
# 	@echo $(call print_error, "Your system not meet any minimum requirements")
# 	@echo $(check_pass)
# 	@exit 1
# endif

# Print start install
start_install:
	@echo "Start installing..."

# Print start install
start_update:
	@echo "Start updating..."

# Build and compile
build:
	$(GO) build -ldflags "-X '$(PACKAGE)/boot.Version=$(VERSION)' -X '$(PACKAGE)/boot.LastBuildAt=$(DATE)'"

# Build and compile under soft maintenance mode
build_soft_maintenance:
	$(GO) build -ldflags "-X '$(PACKAGE)/boot.Version=$(VERSION)' -X '$(PACKAGE)/boot.LastBuildAt=$(DATE)' -X '$(PACKAGE)/boot.SoftMaintenance=true'"

# Build and compile under hard maintenance mode
build_hard_maintenance:
	$(GO) build -ldflags "-X '$(PACKAGE)/boot.Version=$(VERSION)' -X '$(PACKAGE)/boot.LastBuildAt=$(DATE)' -X '$(PACKAGE)/boot.HardMaintenance=true'"


# Linux pre
install_linux: doctor start_install copy_unit_file_linux make_logdir_linux copy_syslog_conf_linux build copy_config mv_to_usr_sbin start_service
update_linux: doctor start_update build mv_to_usr_sbin start_service
smooth_update_linux: doctor start_update build mv_to_usr_sbin
soft_maintenance_linux: doctor start_update build_soft_maintenance mv_to_usr_sbin
hard_maintenance_linux: doctor start_update build_hard_maintenance mv_to_usr_sbin
copy_unit_file_linux:
	cp .dist/linux/$(LINUX_UNIT_FILE) $(shell pwd)/$(LINUX_UNIT_FILE).tmp
	sed -i -e 's/debug_option/$(DEBUG_OPTIONS)/g' $(shell pwd)/$(LINUX_UNIT_FILE).tmp
	mv $(shell pwd)/$(LINUX_UNIT_FILE).tmp $(shell pwd)/$(LINUX_UNIT_FILE)
	sudo cp $(shell pwd)/$(LINUX_UNIT_FILE) /lib/systemd/system
	rm $(shell pwd)/$(LINUX_UNIT_FILE)
	sudo systemctl daemon-reload
	sudo systemctl enable $(SERVICE_IDENTIFIER)
make_logdir_linux:
	sudo mkdir /var/log/$(SERVICE_IDENTIFIER)
	sudo chown syslog:adm /var/log/$(SERVICE_IDENTIFIER)
copy_syslog_conf_linux:
	sudo cp .dist/linux/$(LINUX_SYSLOG_CONF_FILE) /etc/rsyslog.d
	sudo systemctl restart rsyslog.service
	sudo cp .dist/linux/$(LINUX_pId_FILE) /run
copy_config:
	sudo mkdir /etc/$(SERVICE_IDENTIFIER)
	sudo cp  $(shell pwd)/private-key.pem /etc/$(SERVICE_IDENTIFIER)
	sudo cp  $(shell pwd)/public-key.pem /etc/$(SERVICE_IDENTIFIER)
	sudo cp  $(shell pwd)/secret.key /etc/$(SERVICE_IDENTIFIER)
	sudo cp  $(shell pwd)/.config.yml /etc/$(SERVICE_IDENTIFIER)
mv_to_usr_sbin:
	sudo mv $(SERVICE_IDENTIFIER) /usr/sbin
start_service:
	sudo systemctl restart $(LINUX_UNIT_FILE)
	systemctl status $(LINUX_UNIT_FILE)

# Darwin pre
install_darwin: doctor start_install build mv_to_darwin_dir
update_darwin: install_darwin
smooth_update_darwin: install_darwin
mv_to_darwin_dir:
	mv $(SERVICE_IDENTIFIER) .dist/darwin