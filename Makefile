SHELL			:=	/bin/bash
PKGS			:=	subscriber
IMAGES			:=	subscriber
TEST_PKGS		:=	subscriber
FMT_PKGS		:=	subscriber
VERSION			:=	$(shell cat ./VERSION)-$(shell git rev-parse --short HEAD)
LD_FLAGS		:= 	-ldflags "-w -s"


.PHONY: all linux push-images test fmt install deps clean


all: $(addsuffix .out, $(PKGS))

linux: $(addsuffix .linux.amd64, $(PKGS))

push-images: $(addprefix push-image-, $(IMAGES))

images: $(addprefix image-, $(IMAGES))

test: $(addprefix test-, $(TEST_PKGS))

fmt: $(addprefix fmt-, $(FMT_PKGS))

install:  $(addprefix install-, $(PKGS))

deps:
	glide install


clean:
	find . -name "*.out" -type f -delete
	find . -name "*.linux.amd64" -type f -delete


push-image-%:
	docker tag beldpro/$* beldpro/$*:$(VERSION)
	docker push beldpro/$*
	docker push beldpro/$*:$(VERSION)


test-%:
	cd $* && go test ./... -v


install-%:
	cd $* && go install -v


fmt-%:
	cd $* && gofmt -s -w .


image-filla-lb:
	cd infra/docker/lb && docker build -t beldpro/filla-lb .

infra:
	docker-compose \
		-p fillabe \
		-f ./infra/docker/docker-compose.yml \
		up


infra-aws:
	docker-compose \
		-p fillabe \
		-f ./infra/docker/docker-compose.aws.yml \
		up -d


%.pb.go: %.proto
	protoc $^ --go_out=plugins=grpc:.


%.out:
	cd $* && go build $(LD_FLAGS) -v -o $@

%.linux.amd64:
	cd $* && GOOS=linux GOARCH=amd64 GCO_ENABLED=0 go build -a -installsuffix cgo $(LD_FLAGS) -v -o $@
