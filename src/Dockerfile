FROM golang:1.15.6 as go

# Run the code for development
FROM go as run

WORKDIR /src
COPY . .
ENV KUBESTACK_ENV=development
ENTRYPOINT [ "go", "run", "." ]

# Test the code
# This is not used by CI
FROM go as test

WORKDIR /src
RUN go get golang.org/x/lint/golint
COPY . .
ENTRYPOINT ["go", "test", "-v", "./server/", "-race", "-tags=integration"]

# Build binary
FROM go as build

RUN go get github.com/ahmetb/govvv
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -a -o /bin/concurmq .

# Production image
FROM gcr.io/distroless/base as release
COPY --from=build /bin/concurmq /
ENTRYPOINT ["/concurmq"]
