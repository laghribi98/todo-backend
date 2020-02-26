FROM golang:1.13 as builder
WORKDIR ..
RUN mkdir todo-backend
WORKDIR /todo-backend
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-backend ./cmd/todo-backend/main.go

FROM scratch
COPY --from=builder /todo-backend/todo-backend /
ENTRYPOINT ["/todo-backend"]
