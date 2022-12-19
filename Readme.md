## Requirements
To run application docker with docker compose should be installed.

## How to run
To run application execute in cli:
```
cd ./docker-compose
docker-compose up -d
```
## Architecture
The application consists of 3 services:
- object storage for storing csv files (used minio - S3 alternative)
- datasource to access objects from csv file (used clickhouse)
- application itself to serve rest requests

Additionally, object-storage-init used in docker compose to run some initialisation commands for minio server.

Rest endpoint is implemented using Gin library.
Initialisation happens in main.go file. Application starts listening events from minio and will insert data to
clickhouse when file is saved in _promotions_ bucket.

Clickhouse initialised with 2 tables: *promotions* and *promotions_tmp*. Insertion of data is implemented in 3 steps:
- truncate promotions_tmp table (because it has data after previous insertion)
- run data insertion using clickhouse's S3 function [doc](https://clickhouse.com/docs/en/sql-reference/table-functions/s3/#:~:text=Provides%20table%2Dlike%20interface%20to,but%20provides%20S3%2Dspecific%20features.&text=path%20%E2%80%94%20Bucket%20url%20with%20path%20to%20file.)
- run clickhouse's exchange command [doc](https://clickhouse.com/docs/en/sql-reference/statements/exchange/)
Exchange command allows exchange the names of two tables atomically. All these steps are done to comply with requirement
to overwrite storage. Thus, it allows to switch to new data seamlessly.
Notification listener and logic for inserting data is located in **notification_service.go**

**promotions_service.go** contains logic of accessing data by id in clickhouse for rest point.

**promotions_controller.go** contains request handler, it calls **promotions_service.go**
to get data from clickhouse. 

**promotion.go** contains dto that is written as response by controller.

**repository.go** contains logic of connection creation to clickhouse.

After starting application you can access minio ui by http://localhost:9001, e.g. to upload file and make sure
application is working. Credentials to access minio _User_: app, _password_: test12345.
After file uploading it should be inserted into clickhouse and data should be accessible for rest requests.

To test rest endpoint you can use swagger by url http://localhost:8080/swagger/index.html

## Answers for "Additionally, consider" section:
- **The .csv file could be very big (billions of entries) - how would your application perform?**

Based on some articles, clickhouse and s3 function can scale for big amount of data.
To scale we can add resources to our clickhouse cluster. Also, big csv file should be separated to
multiple smaller files to allow parallel processing.
However, in this case, current approach with listening events from S3 bucket will not work.
We can run insert logic by some schedule instead.
Clickhouse S3 function has some settings for parallel processing, e.g. max_threads, etc.
[article 1](https://altinity.com/blog/tips-for-high-performance-clickhouse-clusters-with-s3-object-storage)
[article 2](https://altinity.com/blog/ultra-fast-data-loading-and-testing-in-altinity-cloud)

- **Every new file is immutable, that is, you should erase and write the whole storage;**

This point was considered and described above

- **How would your application perform in peak periods (millions of requests per minute)?**

As example, application can be deployed to Kubernetes services in AWS and scale based on some metrics (cpu, rps)
It can scale when limits of mentioned metrics are exceeded.
By scaling I mean running additional pods to accept load balanced requests.
In this case, we can consider extracting the logic of listening events from S3 to separate service.

- **How would you operate this app in production (e.g. deployment, scaling, monitoring)?**

Scaling covered in previous point.

**Deployment**: we can use some basic tools to deploy application: Jenkins or gitlab-ci.
It should run unit tests, build artifacts, build docker image
and add push it to some docker repository (for further deployment).
To deploy in qa/staging/prod environments we can consider using helm templates
and deploy application using helm upgrade via mentioned ci/cd tools.

**Monitoring**: we can consider using prometheus (for scraping metrics) and grafana to visualise metrics.
Metrics to consider:
1) Rps - requests per seconds, amount of requests application is accepting.
2) Some kubernetes metrics related to memory and CPU consumption in the pod, CPU throttling + state of the pod (if it is alive).
3) Some Go application specific metrics: memory consumption (heap consumption, garbage collection)
4) Minio: we can consider counting of received files.
It can be useful to configure some alert if amount of files is zero for some period of time.
Capacity - how much free space do we have. Perhaps, some disk and I/O metrics.
5) Clickhouse - cpu, memory. Connection pools metrics - idle/active connections.
Some metrics related to partitioning (if any). Metrics related to cluster state (if node is alive).

- **The application should be written in golang;**

Done

- **Main deliverable is the code for the app including usage instructions, ideally in a repo/github gist.**

Done