include make/all.mk

ASDF_VERSION=0.8.1
DOCKER_VERSION=20.10.7

CI_DOCKER_BUILD_ARGS=ASDF_VERSION=$(ASDF_VERSION) \
                     DOCKER_VERSION=$(DOCKER_VERSION)

CI_DOCKER_EXTRA_FILES=.tool-versions .pre-commit-config.yaml
