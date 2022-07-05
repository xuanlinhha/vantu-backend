### SQLite
```sql
create table phrases (
  id int,
  han varchar(100),
  content text,
  info text
)
```

### Build docker images
```sh
docker build -t go-vantu-backend --no-cache=true .
docker tag go-vantu-backend xuanlinhha/go-vantu-backend
docker push xuanlinhha/go-vantu-backend
```

### Local run
```sh
docker run -d --network host --name go-vantu-backend xuanlinhha/go-vantu-backend
```
