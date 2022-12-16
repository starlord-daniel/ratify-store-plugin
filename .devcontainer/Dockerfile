# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.245.2/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.19, 1.18, 1-bullseye, 1.19-bullseye, 1.18-bullseye, 1-buster, 1.19-buster, 1.18-buster
ARG VARIANT="1.19-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 18, 16, 14
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

ARG BATS_VERSION="1.8.2"
RUN curl -Lo bats.tar.gz https://github.com/bats-core/bats-core/archive/v${BATS_VERSION}.tar.gz \
  && tar -zxf bats.tar.gz \
  && bash ./bats-core-${BATS_VERSION}/install.sh /usr/local \
  && rm -rf bats.tar.gz ./bats-core-${BATS_VERSION}
