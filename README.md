# E-commerce system

A learning project: a small e-commerce backend built as microservices in Go. The goal is to practice how real distributed systems are put together — separate services, their own databases, async messaging, and a single API entry point — without pretending to be production-ready.

## What it does

You can register as a customer, log in, browse products, and place orders. An admin manages the catalog and stock. When someone checks out, inventory and notifications react in the background, and the order status moves from `pending` to `confirmed` or `failed`.

## Tech stack

- **Go** + **Gin** — HTTP APIs
- **PostgreSQL** — one database per service
- **Kafka** — async events between services
- **Docker Compose** — run everything locally

## Services

| Service | What it does |
|---------|----------------|
| **gateway** | Single entry point. Routes requests, checks JWTs, blocks non-admins from admin routes. |
| **user-service** | Register, login, `/me`. Issues JWTs. Internal API for other services to look up a user's email. |
| **product-service** | Product catalog — CRUD for admins, public reads. |
| **inventory-service** | Stock per product. Consumes `order.created`, reserves stock, publishes `inventory.reserved` or `inventory.reservation_failed`. |
| **order-service** | Create orders, list/get your orders. Publishes `order.created`, consumes inventory outcome events, updates status. |
| **notification-service** | Consumes inventory outcome events and sends order emails (logged to stdout for now). |

## Order flow (Kafka)

```
POST /orders  →  order.created
                    ↓
              inventory-service  (reserves stock, idempotent per order_id)
                    ↓
         inventory.reserved  OR  inventory.reservation_failed
                    ↓
              order-service  (status → confirmed / failed, idempotent)
                    ↓
              notification-service  (confirmation or failure email)
```

`POST /orders` returns `pending` immediately. Poll `GET /orders/:id` after a few seconds to see the final status.

## API (via gateway)

| Endpoint | Auth | Who |
|----------|------|-----|
| `POST /api/v1/register` | — | anyone |
| `POST /api/v1/login` | — | anyone |
| `GET /api/v1/products` | — | anyone |
| `GET /api/v1/me` | JWT | customer |
| `POST /api/v1/orders` | JWT | customer |
| `GET /api/v1/orders/me` | JWT | customer (list my orders) |
| `GET /api/v1/orders/:id` | JWT | customer (own order only) |
| `POST /api/v1/products` | JWT + admin | admin |
| `PUT /api/v1/inventory/:product_id` | JWT + admin | admin |

## Running it

```bash
docker compose up --build
```

Gateway: `http://localhost:8080`

Default admin: `admin@example.com` / `admin-secret`

Kafka (from host): `localhost:9092`

Prometheus: `http://localhost:9090`

Grafana: `http://localhost:3000` (login `admin` / `admin`) — dashboards for **Gateway**, **Order service**, and **Inventory service**.

## Project layout

Each service is its own Go module (`user-service/`, `product-service/`, etc.). They talk over HTTP or Kafka — no shared Go code between services.
