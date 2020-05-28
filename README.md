# Image media processor
This Application allows a user to resize and upload photos. 



## Main functions
* Application gets image by user request, resize it using user request parameters and upload both images to cloud. In
response, there is an object with links to resized image and original.
* Application allows to see a list of earlier resized images with resizing results and resizing
parameters.
* Application allows to resize old image one more time passing old image id and new resize
parameters.
* API return information about all errors in proper format for invalid user inputs.
* All images returns as an object with links to resized and original image.
* Application has API versioning mechanism.

## Build from source
Building this application requires Go (version 1.14 or later)
```shell
go build
```
Also, if you wish to set name for binary output file, you cat use something like this:
```shell
go build -o image-media-processor
``` 
## Running application
Image media processor application require configuration file to correct work. This file should include configurations about:
* Server settings
* Cloud setting (*at this time - AWS*)
* DB settings (*at this time - MongoDB*)

### Example of ***config.toml*** file
```toml
[Server]
ServerPort = ":8080"
LogLevel = "INFO"

[Aws]
AwsAccessKeyId     = "id"
AwsSecretAccessKey = "accesskey"
AwsRegion          = "eu-central-1"
AwsBucket          = "my-bucket"

[MongoDb]
Username  = "admin"
Password = "admin"
Address = "127.0.0.127017/test"
Store = "imageStore"
Collection = "usersData"
```

## REST Api Examples
| Path | Request example  | 
| :------: | --------- | 
| **/api/v1/resize**    |<pre lang="json">{<br>  "user_id": "a393e097-6f4c-493d-9a82-e612b3d7e53d",<br>  "request_id": "423adce85c",<br>  "width": 100",<br>  "height": 200<br>}</pre>| 
| **/api/v1/resize-by-id**   |<pre lang="json">{<br>  "user_id": "a393e097-6f4c-493d-9a82-e612b3d7e53d",<br>  "request_id": "423adce85c",<br>  "width": 100",<br>  "height": 200,<br>  "image_id": 3652154874<br>}</pre>| 
| **/api/v1/list**    |<pre>user_id=a393e097-6f4c-493d-9a82-e612b3d7e53d&request_id=423adce85c</pre>| 
|

## Docker 
You can build this application using *docker*. This repository include ***Dockerfile*** and ***docker-compose.yaml*** files.
#### Build application using docker
```shell script
docker build -t image-media-processor .
```
#### Run application using docker
```shell script
docker run -p 8080:8080 image-media-processor
```

#### Build application using docker-compose
```shell script
docker-compose build
```
#### Run application using docker-compose
```shell script
docker-compose up -d
```

Be sure, that application config file is included.
