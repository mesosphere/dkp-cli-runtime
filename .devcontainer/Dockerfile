FROM debian:bullseye-20211201-slim

ARG USERNAME=dkp-dev
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt-get update && apt-get install -y \
    bash \
    build-essential \
    bzip2 \
    curl \
    direnv \
    git \
    python3 \
    python3-pip \
    unzip \
    && rm -rf /var/lib/apt/lists/*

RUN pip3 install --upgrade yamale yamllint pre-commit gitlint

ARG DOCKER_VERSION
RUN curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz | \
    tar xz -C /usr/local/bin --strip-components=1 docker/docker

# set up user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME -s /bin/bash
USER $USERNAME
SHELL ["/bin/bash", "-ic"]

# install asdf
ARG ASDF_VERSION
RUN set -eux; \
    git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v$ASDF_VERSION; \
    echo -e ". $HOME/.asdf/asdf.sh" > ~/.bashrc.devcontainer; \
    echo -e "[[ -f ~/.bashrc.devcontainer ]] && source ~/.bashrc.devcontainer" >> ~/.bashrc;
COPY .tool-versions /home/$USERNAME

# install project dependencies
RUN mkdir /tmp/dkp-cli
COPY make /tmp/dkp-cli/make
COPY .tool-versions Makefile /tmp/dkp-cli/
RUN cd /tmp/dkp-cli && git init && make install-tools
