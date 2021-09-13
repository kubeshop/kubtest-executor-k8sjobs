

# TODO add valid dockerfile with executor
FROM golang:1.17
# RUN apk update -y 
RUN apt update && apt install -y nodejs npm
RUN npm install -g newman
WORKDIR /agent
COPY pkg.newman/agent /agent/agent
# RUN go build -o agent

# FROM postman/newman
# RUN apk --no-cache add ca-certificates
