FROM golang:1.8

ENV SRC_FOLDER /go/src/payment-service
ENV CONFIG_FOLDER /go/config
ENV PKG_FOLDER /go/pkg
ENV MONGO_URL mongodb://mongo:27017/payment-db
RUN mkdir -p $SRC_FOLDER $PKG_FOLDER

WORKDIR $SRC_FOLDER

COPY . $SRC_FOLDER

RUN go install && rm -rf $PKG_FOLDER

HEALTHCHECK --interval=15s --retries=10 CMD curl -fs http://localhost:8080/health || exit 1

EXPOSE 8080

CMD /go/bin/payment-service