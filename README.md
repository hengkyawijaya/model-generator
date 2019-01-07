### How To Use

* Define your schema and configuration in json folder
* Install the dependencies
  ``` 
  dep ensure -v -vendor-only
  ```
* Run the generator  
  ```
  go run generator.go
  ```
* You can see the magic in the folder migration
* Build docker image postgres 
  ```
  docker-compose build
  ```
* Start the container postgres
  ```
  docker-compose up
  ```
