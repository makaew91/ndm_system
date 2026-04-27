FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /broker .

FROM scratch
COPY --from=build /broker /broker
EXPOSE 8080
ENTRYPOINT ["/broker"]
CMD ["8080"]
