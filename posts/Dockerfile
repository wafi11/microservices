FROM rust:1.88-bookworm AS builder
WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs \
    && cargo build --release

COPY ./src ./src
RUN touch src/main.rs && cargo build --release

FROM debian:bookworm-slim
WORKDIR /usr/local/bin

RUN apt-get update && apt-get install -y \
    libssl3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/app/target/release/posts_service ./posts_service
CMD ["./posts_service"]