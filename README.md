#Image Previewer

## Service Description
The service is intended for making a preview (creating a small copy from an image)

## Commands for working with the project
1. make build-server - build application binary
2. make test-single - run one test pass with 1m timeout
3. make test-race - runs 100 test passes with a timeout of 7m
4. make test-coverage - see % of project test coverage
5. make test-integration - run integration tests
6. make lint - run linter in the project

### Run docker
Docker is started with the `--remove-orphans` flag to avoid cluttering up old containers

1. make start - starts the project
2. make stop - stop containers

## Description of the service
We send parameters to url resize to resize the image and a link to get the original image we want to change

### url parameters
`host/fill/{width}/{height}/{imageUrl}`

1. width - the width of the image we want to get
2. height - the height of the image we want to get
3. imageUrl - link to the original image to reduce by the specified parameters

### Example for getting an image after running docker
`http://127.0.0.1/fill/200/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_1024x252.jpg`