FROM golang:1.14-alpine AS build

RUN apk update && apk upgrade && apk add --no-cache bash git

WORKDIR .
COPY . /proj/

RUN /bin/bash
RUN cd /proj && CGO_ENABLED=0 go build github.com/simonmittag/mse6/cmd/mse6

#multistage build uses output from previous image
FROM alpine
COPY --from=build /proj/mse6 /mse6
ENV wait 3
EXPOSE 8081
ENTRYPOINT "/mse6" "-w=${wait}"