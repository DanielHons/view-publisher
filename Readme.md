Create a view that returns a json mode (content type will be set to application/json) in PostgreSQL.

# Build
```
go build .
```

#Run
```
create view my_custom_report as 
....
```
Configure the name of the view for the desired subpath in a config file
```
conn: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
port: 3333
publishers:
  desired-sub-path: my_custom_report
```


Run the app with the configuration
```
./view-publisher -c /path/to/config.yaml
```

Now you should find your json calling
```
GET http://localhost:3333/desired-sub-path
```
