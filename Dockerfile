FROM golang:latest as builder

ADD . /go/src/sfladmin
WORKDIR /go/src/sfladmin
RUN go get github.com/futurenda/google-auth-id-token-verifier
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/lib/pq
RUN go get github.com/joho/godotenv
RUN go get github.com/jinzhu/gorm
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go install sfladmin
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/sfladmin/main .
COPY --from=builder /go/src/sfladmin/.env .
COPY --from=builder /go/src/sfladmin/pgserver.crt .
COPY --from=builder /go/src/sfladmin/pgserver.pem .

# Expose port 6060 to the outside world
EXPOSE 6060

# Command to run the executable
CMD ["./main"] 