# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.202.5/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && curl -s https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash \
    && apt-get -y install --no-install-recommends cmake git-lfs \
    && git clone https://github.com/bab2min/Kiwi.git \
    && cd Kiwi \
    && git lfs pull \
    && git submodule sync \
    && git submodule update --init --recursive \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release ../ \
    && make && make install \
    && ldconfig \
    && cd ../.. && rm -rf Kiwi \
    && apt-get --purge -y remove cmake git-lfs \
    && apt autoremove \
    && rm -rf /var/lib/apt/lists/*


RUN sh -c "$(curl -fsSL https://starship.rs/install.sh)" -- --yes

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1