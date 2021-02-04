FROM golang

WORKDIR /src
COPY . .

RUN go mod tidy
# RUN ./unit_test.sh

RUN go build -o main
CMD ["./main"]