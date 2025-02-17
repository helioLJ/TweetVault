FROM golang:1.23-alpine

WORKDIR /app

# Install air for hot reloading (using a compatible version)
RUN go install github.com/cosmtrek/air@v1.49.0

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Expose port 8080
EXPOSE 8080

# Run air for development
CMD ["air"] 