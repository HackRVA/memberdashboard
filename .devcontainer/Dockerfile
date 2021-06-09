# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.177.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT="1.15"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Option] Install Node.js
ARG INSTALL_NODE="true"
ARG NODE_VERSION="12.15"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

RUN curl -o /usr/local/bin/swagger -L'#' "https://github.com/go-swagger/go-swagger/releases/download/v0.27.0/swagger_linux_amd64" \
    && chmod +x /usr/local/bin/swagger \
    && git clone https://github.com/vishnubob/wait-for-it.git \
    && mv wait-for-it/wait-for-it.sh /usr/local/bin/ \
    && chmod +x /usr/local/bin/wait-for-it.sh

RUN apt update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends postgresql

RUN go get -x -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

WORKDIR /workspace
# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

