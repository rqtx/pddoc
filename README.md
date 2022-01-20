# pddoc

A tool to generate documentation from cloud providers in markdown format

Docker run example

``` bash
docker run --rm -v $(pwd):/doc -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY docker.com/rqtx/pdoc:0.0.0- aws -f /doc/test.md -r sa-east-1
```
