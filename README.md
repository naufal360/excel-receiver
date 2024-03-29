# Excel Receiver

Upload media file for format .xlsx and .csv and save data from file into the Artemis queue, insert into MySQL and save the file in local storage directory.

# Docs
- [Flowchart](https://drive.google.com/file/d/1WqvIM0Nae-JRj61W3ALtSUFvzwPhkiHI/view?usp=sharing)
- [SonarQube](http://192.168.181.116:9090/dashboard?id=NaufalSimpleFileAPI)

## Requirements
- Golang 1.19.6 for (development)
- Docker
- MySQL
- Artemis


## Development
```bash
cp config.yaml.example config.yaml
go mod tidy
go run main.go
```

## Unit Test
### Unit test only
```bash
go test ./... -v
``` 
### Unit test for SonarQube coverage.
- This will create coverage.out and report.json files that will be included inside sonar.properties
```bash
go test "./..." -coverprofile="coverage.out" -covermode=count -json > report.json;
```

## Sonarqube Scan (Manual)
- scan using docker
```bash
docker run \
    --rm \
    -v "$(pwd):/usr/src" \
    -v "$(pwd)/sonar.properties:/opt/sonar-scanner/conf/sonar-scanner.properties" \
    sonarsource/sonar-scanner-cli:4.7
```
## Deployment (Docker Build)
- Create Database on your MySQL, command line refer to file [/initdb/init.sql](/initdb/init.sql)
- Copy (create) config file refer to `config.yaml.example`
```bash
cp config.yaml.example config.yaml
```
- Update config if required such as (db connection)
- if you are updating the port config, make sure you update the port in dockerfile or use `ENV PORT=`
- Build the app using docker
```bash
docker build -t excel-receiver:1.0.0 .
```
- Run the images below, then the app is ready to use.
```bash
docker run -d -p 5000:5000 --name excel-receiver excel-receiver:1.0.0 
```

