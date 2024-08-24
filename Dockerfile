# officle golang image
FROM golang:1.22.0
LABEL maintainers="Malak and Mariam and Mariam"
LABEL version="1.0"
LABEL description="web server written in Go to generate ASCII art."
# set work directory
WORKDIR /app
  
#copy the mod and server qand everything else needed to run the app
COPY . .
# Build the Go app
RUN go build -o main .
# expose the port
EXPOSE 5500
# run the app
CMD ["/app/main"]

# docker build -t my-go-app .
# docker run -d -p 5500:5500 my-go-app