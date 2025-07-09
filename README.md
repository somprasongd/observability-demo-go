# Observability

> Code example for Observability (Logs, Traces, Metrics) with Go, Framework Fiber, tools OpenTelemetry, and popular stacks Grafana, Loki, Tempo, Prometheus

---

## Why Observability is important?

Today's software systems run on multiple machines or servers, we have:

- Distributed Systems
- Microservices
- Dynamic Cloud Infrastructure

**Question:** If a Service fails or data is incorrect, how do you find the root cause?

Without Observability, you'd have to `ssh` into the machine and manually check logs or debug code, which is no longer efficient.

---

## Monitoring vs Observability

First, let's differentiate between **Monitoring** and **Observability**

| Monitoring | Observability |
| --- | --- |
| Sets alerts for abnormal conditions (e.g., high CPU) | Answers the question "Why" a condition occurred |
| Only knows what happened | Knows where it happened, how it happened, and where to fix it |

Example: A Dashboard shows error rate | Example: Traces and Logs are connected

**Monitoring = What**, **Observability = Why**

---

## Three Components of Observability

For good Observability, we need 3 things:

- **Logs** — Event logs (e.g., errors, debugs, info)
- **Metrics** — Numerical summaries of state (e.g., requests/second, CPU, memory)
- **Traces** — The journey of a request across Services

With all three, you'll get:

- A high-level overview
- Detailed insights
- Faster root cause detection
- Quicker bug fixing

---

## Why use OpenTelemetry + Grafana Stack?

**OpenTelemetry (OTel)** is the open standard for storing Logs, Metrics, Traces in a unified format.

**Grafana** is the popular visualization tool.

- **Loki** stores Logs
- **Tempo** stores Traces
- **Prometheus** or **Mimir** stores Metrics

All integrate well, are easy to use, have a large community, and are free!

---

## Summary

- **Observability** helps you understand "Why" a system is malfunctioning.
- Covers Logs + Metrics + Traces.
- We'll be using Go + Fiber + OpenTelemetry + Loki + Tempo + Prometheus + Grafana.

---
