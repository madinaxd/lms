# LMS

Currently App contains of two microservices Students and Courses. 

### Microservices.

Microservices are written using go-kit tool.
* Go-kit is designed to be modular and composable, which means you can easily plug in and swap out different components such as service discovery, load balancing, and transport layers. This makes it easier to create microservices that are flexible and scalable. 
* Go-kit is designed to work well with other tools in the microservices ecosystem, such as Kubernetes, Consul, and Prometheus. This makes it easier to deploy and manage your microservices, as well as to monitor and debug them.

### DB

* PostgreSQL is used as a main database.It has strong data integrity features, which ensures that data stored in the database is consistent and accurate. This is critical for an LMS, which must maintain accurate records of student progress and performance.

* It is highly scalable and can handle large amounts of data and high volumes of traffic. This is important for an LMS, which may have many users accessing the system simultaneously.

* PostgreSQL has a large and active community of developers

### Monitoring 

Prometheus is used to create metrics

### Logging 

Logging is done with "go-kit/log" package, which allows to use different logging backends such as stdout, files, and remote services, among others. Also go-kit logging interface defines a set of methods for different logging levels, such as Info, Error, and Debug

## Local Development

Run the commands below.

```console
docker-compose -up -d
```

## Test

Now you can test the API with Postman at:

```console
localhost:8081/students
localhost:7071/courses
```


## Further development

To create API Gateway. 

```