FROM golang:1.20.4-alpine AS build_stage
COPY ./ /go/src/app/
WORKDIR /go/src/app
RUN  go mod download
RUN  go install ./cmd/http/main.go

FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/bin/main /app_binary/
RUN chmod +x main
EXPOSE 8081/tcp
ENTRYPOINT ./main

EXPOSE 8080/tcp
CMD [ "main" ]

# docker build --file=build/package/http/Dockerfile --tag=httplinks:latest .
