MICROSERVICE:=kc_rslbe
GO:=go
GOFLAGS:=-v -p 1

ifeq ($(GOARCH),arm64)
DOCKER_TAG=${GOARCH}
GOARCH=arm64
else
GOARCH=amd64
DOCKER_TAG=latest
endif

KC_COMMON_JAR_VERSION:=$(shell bin/extract_kc_common_version_from_pom.sh kc_rslbe_dl3/pom.xml)
KC_COMMON_JAR:=${HOME}/.m2/repository/com/kaiostech/kc_common/${KC_COMMON_JAR_VERSION}/kc_common-${KC_COMMON_JAR_VERSION}.jar

CLOUD_TARGETS:=bin/${MICROSERVICE}_ll3 bin/${MICROSERVICE}_fe3
TOOL_TARGETS:=bin/gen_doc

VERSION_DOCKER=$(shell echo -e "{\n    \"${MICROSERVICE}\" : {\n    \"docker_tag\" : \"${DOCKER_TAG}\",\n    \"arch\" : \"${GOARCH}\",\n    \"version\" : \"$$(cat .version)\",\n    \"tag\":\""$$(cat .version)"\"\n  }\n}" > bom.json)
GITREF:=$(shell if [ -d .git ]; then git describe --tags --abbrev=9; else echo "none"; fi)

default: githook ${MICROSERVICE}_dl3 ${CLOUD_TARGETS} ${TOOL_TARGETS}

.PHONY: ${MICROSERVICE}_dl3 ${MICROSERVICE}_ll3 ${MICROSERVICE}_fe3 utils

all: githook clean default

############# Hook start
hook: githook

githook: .git/hooks/pre-commit

.git/hooks/pre-commit: bin/pre-commit
	@thedate=`date +"%Y%m%d_%H%M%S"`; if [ -f .git/hooks/pre-commit ]; then echo "Updating Git pre-commit hook on $${thedate}" ; cp .git/hooks/pre-commit .git/hooks/pre-commit.backup.$${thedate}; else echo "Installing Git pre-commit hook"; fi; mkdir -p .git/hooks && touch .git/hooks/pre-commit && cp -f bin/pre-commit .git/hooks/pre-commit && chmod 755 .git/hooks/pre-commit

format:
	mvn -X -e -f ${MICROSERVICE}_dl3/pom.xml git-code-format:format-code -Dgcf.globPattern=**/*
	@for f in `find . -name \*.go`; do echo "Reformatting $$f"; go fmt $$f || { echo "Issue found in Golang code. Please fix them and try again."; exit 1; }; done

vet:
	@GO_MODULE_NAME=`head -1 go.mod | awk '{ print $$2 }'`; for path_to_test in `find . -type d -print`; do  { ls >/dev/null 2>&1 -al $${path_to_test}/*.go || continue ; };  path_to_test=`echo $${path_to_test} | sed 's/^.\///g'`; go_path=$${GO_MODULE_NAME}/$${path_to_test}; echo -n "Vetting $${go_path} ... "; ${GO} vet $${go_path}; if [ $$? -ne 0 ]; then echo \"Issue found in Golang code. Please fix them and try again.\"; exit 2; else echo "[ OK ]"; fi done
############# Hook end

${MICROSERVICE}_fe3: githook bin/${MICROSERVICE}_fe3

${MICROSERVICE}_ll3: githook bin/${MICROSERVICE}_ll3

${MICROSERVICE}_dl3: githook ${KC_COMMON_JAR} kc_rslbe_dl3/pom.xml
	sh -c 'cd ${MICROSERVICE}_dl3; make'

${MICROSERVICE}_ll3/version/version.go: .version
	sh -c 'mkdir -p ${MICROSERVICE}_ll3/version && ./bin/gen_version_go.sh .version >$@'

${MICROSERVICE}_fe3/version/version.go: .version
	sh -c 'mkdir -p ${MICROSERVICE}_fe3/version && ./bin/gen_version_go.sh .version >$@'

bin/${MICROSERVICE}_ll3: ${MICROSERVICE}_ll3/version/version.go
	@echo "========== Compiling $@ =========="
	sh -c '$(GO) build $(GOFLAGS) -o $@ git.kaiostech.com/cloud/${MICROSERVICE}/${MICROSERVICE}_ll3'

bin/${MICROSERVICE}_fe3: ${MICROSERVICE}_fe3/version/version.go
	@echo "========== Compiling $@ =========="
	sh -c '$(GO) build $(GOFLAGS) -o $@ git.kaiostech.com/cloud/${MICROSERVICE}/${MICROSERVICE}_fe3'

utils/version/version.go:
	sh -c 'mkdir -p utils/version && ./bin/gen_version_go.sh .version >$@'


${KC_COMMON_JAR}: kc_rslbe_dl3/pom.xml
	sh -c "./bin/gen_kc_common.sh ${KC_COMMON_JAR_VERSION}"


bin/gen_doc: utils/version/version.go
	@echo "========== Compiling $@ =========="
	@sh -c 'binary_name=`basename $@`;  echo $(GO) build $(GOFLAGS) -o $@ git.kaiostech.com/cloud/${MICROSERVICE}/utils/$${binary_name} && $(GO) build $(GOFLAGS) -o $@ git.kaiostech.com/cloud/${MICROSERVICE}/utils/$${binary_name}'


deploy: clean ${MICROSERVICE}_dl3 ${CLOUD_TARGETS} ${KC_COMMON_JAR}
	version=`cat .version`; output_folder=${MICROSERVICE}-$${version}; mkdir -p /tmp/$${output_folder}/bin && cp -a bin/{${MICROSERVICE}_fe3,${MICROSERVICE}_ll3,run3.sh,service_functions.sh,start_service.sh,stop_service.sh,${MICROSERVICE}_dl3-$${version}.jar,functions.sh,jars,conf,run_${MICROSERVICE}_dl3.sh} /tmp/$${output_folder}/bin/ && cp ${KC_COMMON_JAR} /tmp/$${output_folder}/bin/jars/ && mkdir -p /tmp/$${output_folder}/conf/ && cp -a RELEASE.txt /tmp/$${output_folder} && cp -a conf/{${MICROSERVICE}.conf.dev,${MICROSERVICE}.conf.local,${MICROSERVICE}.conf.k5,jwt_key.priv,jwt_key.pub,${MICROSERVICE}_fe3_log.xml,${MICROSERVICE}_ll3_log.xml} /tmp/$${output_folder}/conf/ && tar -jcvf /tmp/${MICROSERVICE}-$${version}-${GOARCH}.tar.bz2 -C /tmp ${MICROSERVICE}-$${version} && mv /tmp/${MICROSERVICE}-$${version}-${GOARCH}.tar.bz2 . && rm -Rf /tmp/${MICROSERVICE}-$${version}

docker: ${VERSION_DOCKER}
	@sh -c 'tiller_json=`cat bom.json` tiller -b . -n'
	docker build . -t ${MICROSERVICE}-${GOARCH}:${GITREF}

ci-container: ${VERSION_DOCKER}
	@sh -c 'tiller_json=`cat bom.json` tiller -b . -n'

doc: bin/gen_doc
	sh -c './bin/gen_doc'
	sh -c 'raml2html kaicloud.raml > ${MICROSERVICE}.html'

test:
	@echo "Testing ${MICROSERVICE}_fe3 ..." && find ${MICROSERVICE}_fe3 -type d -exec go test git.kaiostech.com/cloud/${MICROSERVICE}/{} \;
	@echo "Testing ${MICROSERVICE}_ll3 ..." && find ${MICROSERVICE}_ll3 -type d -exec go test git.kaiostech.com/cloud/${MICROSERVICE}/{} \;
	@echo "Testing ${MICROSERVICE}_dl3 ..." && cd ${MICROSERVICE}_dl3 && make test

clean:
	cd ${MICROSERVICE}_dl3 && make clean
	@echo "Deleting generated binary files ..."; for binary in ${CLOUD_TARGETS} ${TOOL_TARGETS} ; do if [ -f "$${binary}" ]; then rm -f $${binary} && echo $${binary}; fi; done
	@echo "Deleting generated version files ..."; for version_dir in ${MICROSERVICE}_ll3/version ${MICROSERVICE}_fe3/version utils/version; do if [ -d "$${version_dir}" ]; then rm -Rf $${version_dir} && echo $${version_dir}; fi; done
	@echo "Deleting generated bom.json file ..." && rm -f bom.json
	@echo "Deleting emacs backup files ..."; find . -type f -name \*~ -exec rm {} \; -print
	@echo "Deleting log files ..."; find . -maxdepth 1 -type f \( -name \*.log.\* -o -name \*.log \) -exec rm {} \; -print
	@echo "Deleting Version information ..."; if [ -d version ]; then rm -Rf version && echo version; fi
	@generated_archive=`ls ${MICROSERVICE}-*${GOARCH}.tar.bz2 2>/dev/null`; if [ "x" != "x$${generated_archive}" ] && [ -f "$${generated_archive}" ]; then echo "Deleting $${generated_archive} ..."; rm -f $${generated_archive}; fi

clear_cache:
	#@echo "Deleting local caches ..."; for cache in ~/go/pkg/mod/cache/download/git.kaiostech.com/cloud/common ~/go/pkg/mod/cache/download/git.kaiostech.com/cloud/thirdparty ~/go/pkg/mod/git.kaiostech.com/cloud ; do if [ -d "$${cache}" ]; then find $${cache} -type f -exec chmod 644 {} \; -print ; find $${cache} -type d -exec chmod 755 {} \; -print ; rm -Rf $${cache} && echo $${cache}; fi; done
	find ~/go -type d -exec chmod 755 {} \; && rm -Rf ~/go