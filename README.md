# Babel Bot

For Babel Street Interview

## Development + Deployment

To spin up a container with hot reloading, make your own file *.env* and place it the root of this project. Put in it these variables:
```
MYSQL_ROOT_PASSWORD=<some root password>
MYSQL_DATABASE=<some name>
MYSQL_USER=<some user>
MYSQL_PASSWORD=<some password>

TWITTER_CONSUMER_KEY=<consumer key genned by twitter>
TWITTER_CONSUMER_SECRET=<consumer secret genned by twitter>
TWITTER_ACCESS_TOKEN=<access token genned by twitter>
TWITTER_ACCESS_SECRET=<access secret genned by twitter>
```

To actually run the app:
```
docker-compose up
```