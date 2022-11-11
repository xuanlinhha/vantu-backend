# Build
FROM ubuntu
WORKDIR /app
COPY . .
EXPOSE 3000
ENTRYPOINT ["/app/vantu-backend"]
