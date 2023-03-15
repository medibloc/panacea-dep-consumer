FROM golang:1.19.2-bullseye AS build-env

# Install minimum necessary dependencies,
ENV PACKAGES make git gcc
RUN apt-get update -y
RUN apt-get install -y $PACKAGES

COPY . /src/panacea-dep-consumer

WORKDIR /src/panacea-dep-consumer

RUN make clean && make build

FROM debian:bullseye-slim

COPY --from=build-env /src/panacea-dep-consumer/build/consumerd /usr/bin/consumerd

RUN chmod +x /usr/bin/consumerd

CMD ["/consumerd"]

EXPOSE 8080
