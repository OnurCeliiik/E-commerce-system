# E-commerce system

A learning project: a small e-commerce backend built as microservices in Go. The goal is to practice how real distributed systems are put together — separate services, their own databases, async messaging, and a single API entry point — without pretending to be production-ready.

## What it does

You can register as a customer, log in, browse products, and place orders. An admin can manage the product catalog and stock levels. When someone checks out, other services react in the background (inventory updates, notification emails).

## Tech stack

- **Go** + **Gin** — HTTP APIs
- **PostgreSQL** — one database per service
- **Kafka** — events between services (e.g. when an order is created)
- **Docker Compose** — run everything locally

## Services

| Service | What it does |
|---------|----------------|
| **gateway** | Front door. Routes requests, checks JWTs, blocks non-admins from admin routes. |
| **user-service** | Register, login, `/me`. Issues JWT tokens. |
| **product-service** | Product catalog — create, list, update, delete products. |
| **inventory-service** | Stock per product. Listens to `order.created` and subtracts quantity. |
| **order-service** | Place orders. Looks up prices from product-service, saves the order, publishes `order.created`. |
| **notification-service** | Listens to `order.created` and sends a confirmation email (logged to stdout for now). |

## Running it

```bash
docker compose up --build
```

Gateway: `http://localhost:8080`

Default admin (from compose env): `admin@example.com` / `admin-secret`

Kafka (from host): `localhost:9092`

## Project layout

Each service is its own Go module under its own folder (`user-service/`, `product-service/`, etc.). They talk over HTTP or Kafka — no shared code between services.

This is intentionally a work in progress. More pieces (monitoring, real email, tests, etc.) will come later as the project grows.
