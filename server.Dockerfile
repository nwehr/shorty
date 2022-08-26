FROM golang:1.18
WORKDIR /app
COPY . .
RUN make server
CMD /app/server