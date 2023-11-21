FROM go:1.21.1

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o warpbuild-agentd cmd/agentd/main.go

CMD ["./warpbuild-agentd"]