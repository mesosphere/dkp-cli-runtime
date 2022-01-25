# Copyright 2022 D2iQ, Inc. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

OS := $(shell uname -s)
ifeq ($(OS), Darwin)
  BREW_PREFIX := $(shell brew --prefix)
  ifeq ($(BREW_PREFIX),)
    $(error Unable to discover brew prefix - do you have brew installed? See https://brew.sh/ for details of how to install)
  endif

  GNUBIN_PATH := $(BREW_PREFIX)/opt/coreutils/libexec/gnubin
  ifeq ($(wildcard $(GNUBIN_PATH)/*),)
    $(error Cannot find GNU coreutils - have you installed them via brew? See https://formulae.brew.sh/formula/coreutils for details)
  endif
  export PATH := $(GNUBIN_PATH):$(PATH)

  GNUFIND_PATH := $(BREW_PREFIX)/opt/findutils/libexec/gnubin
  ifeq ($(wildcard $(GNUFIND_PATH)/*),)
    $(error Cannot find GNU findutils - have you installed them via brew? See https://formulae.brew.sh/formula/findutils for details)
  endif
  export PATH := $(GNUFIND_PATH):$(PATH)
endif
