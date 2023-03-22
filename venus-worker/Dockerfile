ARG RUST_TOOLCHAIN

FROM rust:${RUST_TOOLCHAIN} AS build

RUN apt-get update
RUN apt-get install -y \
    --no-install-recommends \
    make \
    mesa-opencl-icd \
    ocl-icd-opencl-dev \
    libhwloc-dev \
    libssl-dev \
    pkg-config \
    libclang-dev

WORKDIR /build/
COPY . /build/

ARG GIT_FALLBACK="Unknown (no git or not git repo)"
ARG GIT_COMMIT=${GIT_FALLBACK}
RUN make build-all

FROM debian:bullseye
RUN apt-get update && \
    apt-get install -y \
    --no-install-recommends \
    libhwloc-dev \
    libnuma-dev \
    mesa-opencl-icd \
    ocl-icd-opencl-dev

# DIR for app
WORKDIR /app
ENV PATH "$PATH:/app/"

COPY --from=build /build/target/release/venus-worker /app/venus-worker
COPY --from=build /build/assets/venus-woker.default.toml /app/config.toml

VOLUME [ "/data" ]

EXPOSE 17891
ENTRYPOINT ["venus-worker"]
