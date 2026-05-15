# Stage 1: build Vite client
FROM node:20-alpine AS client-build
WORKDIR /app/client
COPY client/package*.json ./
RUN npm ci
COPY client/ ./
RUN npm run build

# Stage 2: build Go server
FROM golang:1.26-alpine AS server-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.version=${VERSION}" \
    -o hexar-server \
    ./cmd/server

# Stage 3: minimal runtime (no shell, no package manager)
FROM gcr.io/distroless/static-debian12
COPY --from=server-build /app/hexar-server /hexar-server
COPY --from=client-build /app/client/dist /client/dist
EXPOSE 8080
ENTRYPOINT ["/hexar-server", "-addr", ":8080", "-client", "/client/dist"]
