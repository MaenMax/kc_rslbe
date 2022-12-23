GOPATH:=$(shell pwd)
MICROSERVICE:=kc_rslbe
MICROSERVICE_BIN:=kc_rslbe
GO:=go 
GOFLAGS:=-v -p 1


CLOUD_TARGETS:= bin/${MICROSERVICE_BIN}_fe3 


default: ${CLOUD_TARGETS}


bin/${MICROSERVICE_BIN}_fe3: 
	@echo "========== Compiling $@ =========="
	sh -c '$(GO) build $(GOFLAGS) -o $@ git.kaiostech.com/cloud/${MICROSERVICE}/${MICROSERVICE}_fe3'


clean:
	@echo "Deleting generated binary files ..."; for binary in ${CLOUD_TARGETS}; do if [ -f "$${binary}" ]; then rm -f $${binary} && echo $${binary}; fi; done
