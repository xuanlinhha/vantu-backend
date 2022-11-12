# Build
FROM golang:1.18-buster AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o api-server vantu.org/go-backend/cmd/api-server

# Deploy
FROM ubuntu
WORKDIR /app
COPY --from=build /app/api-server ./
COPY data .
# RUN echo "$(ls data)"
# RUN echo "$(pwd)"
EXPOSE 3000

ENTRYPOINT ["/app/api-server"]
