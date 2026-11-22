FROM node:22-bookworm-slim

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN mkdir -p /home/node/.local/share \
    && chown -R ${USER_ID}:${GROUP_ID} /home/node

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    git \
    git-lfs \
    jq \
    less \
    make \
    openssh-client \
    ripgrep \
    unzip \
    && rm -rf /var/lib/apt/lists/*

RUN npm install -g @opencode/cli

WORKDIR /workspace

USER node

CMD ["sleep", "infinity"]
