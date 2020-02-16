FROM golang:1.13-alpine as builder
WORKDIR ..
RUN mkdir todo-backend
WORKDIR /todo-backend
COPY ../.. ./
RUN go build -o todo-backend ./cmd/main.go

FROM scratch
COPY --from=builder /todo-backend/todo-backend .
CMD ["./todo-backend"]
