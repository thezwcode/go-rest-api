# go-rest-api
Attempt to build RestAPI

#API Setup

This application was build with Mac OS Catalina 10.1.15 and assumes that [Golang](https://golang.org/doc/install)v1.13.4 and [PostgreSQL](https://www.postgresql.org/download/)v12.1 are already set up. For more information on setting them up, click on the respective links.
[VSCode](https://code.visualstudio.com/download) was used as the code editor, with Go extension installed.

In the command line, use command ```go env GOPATH``` to check $GOPATH directory. For my machine, $GOPATH=/Users/{username}/go.

If you want to develop the project in another directory, create the directory, navigate there and use command ```export GOPATH=$PWD``` to change $GOPATH to that directory.

In $GOPATH directory, enter:
```
go get -u github.com/gorilla/mux
go get -u github.com/lib/pq
go get -u github.com/joho/godotenv
go get -u github.com/golang/gddo/httputil/header

```

##Database configuration

First, we need to configure the database setup. Ensure that PostgreSQL server is running. 

Using a superuser is not recommended. A new dbuser can be created the command line with this command:

```
createuser {dbuser} --pwprompt --createdb
```

You will be prompted to input your desired password:

```
Enter password for new role: {dbuser password}
Enter it again: {dbuser password}
```

Next, go into the postgresql command prompt:
```
psql postgres -U {dbuser}
```
Create new database:
```
postgres=> CREATE DATABASE {dbname};
```
Grant all privileges for the new database to {dbuser} and see the list of users:
```
postgres=> GRANT ALL PRIVILEGES ON DATABASE {dbname} TO {dbuser}; \list
```
Connect to the database:
```
postgres=> \connect {dbname}
```

Exit:
```
\q
```
##Test database connection

Create a .env file in $GOPATH directory with these attributes (port number is default PostgreSQL port number):
'''
APP_DB_USERNAME = "{dbuser}"
APP_DB_PASSWORD = "{dbuser password}"
APP_DB_NAME = "{dbname}"

SSL_MODE="disable"
PORT="5432"
'''
Log into PostgreSQL: 
```
psql postgres -U {dbuser}
\connect {dbname}
```
##Running the server
Navigate to file directory with go files inside:
```
cd $GOPATH
cd go-rest-api
```
Enter this command in the command line to run the server:
go run *.go

##PostMan

[PostMan](https://www.getpostman.com/downloads/) is used for testing the API, by sending requests and receiving responses.

After setting up PostMan, create a new request. Set the request URL to be ```http://localhost:8080```
The following requests are valid:
```
GET: /search/{objectType}/{id}
POST: /update/{objectType}, send JSON body with attributes other than id
PUT: /update/{objectType}/{id}, send JSON body with attributes other than id
```
To test requests without JSON, let the header remain as any option other than JSON. To change to JSON, go to "Headers", click on "Presets", then "Manage Presets".

Add "application/json".

###JSON strings

To send requests in JSON strings, go to "Body" tab and input JSON object. Only GET, POST and PUT requests work, with their specified endpoints, as mentioned above.

Thanks for trying the program out and feel free to give any feedback (:
