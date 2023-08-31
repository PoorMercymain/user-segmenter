FROM golang:1.20
WORKDIR /user-segmenter
COPY go.mod go.sum ./
RUN go mod download
COPY . /user-segmenter
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-segmenter/cmd/bin/main /user-segmenter/cmd/main.go
CMD ["bash", "-c", "/user-segmenter/cmd/bin/main"]