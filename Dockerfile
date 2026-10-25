FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
COPY public ./public
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/basbas .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/basbas /basbas
ENV PORT=8080
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/basbas"]
