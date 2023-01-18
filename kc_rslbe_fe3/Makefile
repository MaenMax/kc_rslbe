GOPATH:=$(shell pwd)
MICROSERVICE:=kc_rslbe
MICROSERVICE_BIN:=kc_rslbe
GO:=go 
GOFLAGS:=-v -p 1


KC_COMMON_JAR_VERSION:=$(shell bin/extract_kc_common_version_from_pom.sh kc_rslbe_dl3/pom.xml)
KC_COMMON_JAR:=${HOME}/.m2/repository/com/kaiostech/kc_common/${KC_COMMON_JAR_VERSION}/kc_common-${KC_COMMON_JAR_VERSION}.jar


CLOUD_TARGETS:= bin/${MICROSERVICE_BIN}_fe3 


default: ${CLOUD_TARGETS}


bin/${MICROSERVICE_BIN}_fe3: 
	@echo "========== Compiling $@ =========="
	sh -c '$(GO) build $(GOFLAGS) -o $@ git.kaiostech.com/cloud/${MICROSERVICE}/${MICROSERVICE}_fe3'


clean:
	@echo "Deleting generated binary files ..."; for binary in ${CLOUD_TARGETS}; do if [ -f "$${binary}" ]; then rm -f $${binary} && echo $${binary}; fi; done
