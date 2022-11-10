# Build
FROM golang:1.18-buster AS build
WORKDIR /app
COPY . .
# RUN echo "$(ls)"
RUN go mod tidy
RUN go build -o ./go-vantu-backend vantu.org/go-backend/cmd/vantu-backend

# Deploy
FROM ubuntu:21.10
WORKDIR /app
COPY --from=build /app/go-vantu-backend /app/.env /app/vantu.db ./
COPY --from=build /app/static_data/ ./static_data/
# RUN echo "$(ls)"
# RUN echo "$(pwd)"
EXPOSE 3000
ENTRYPOINT ["/app/go-vantu-backend"]
