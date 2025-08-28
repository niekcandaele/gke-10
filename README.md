# GKE-10 Hackathon: AI-Powered Economic Simulation

## 🎯 Project Overview

This repository implements an **AI-powered economic simulation ecosystem** for the [GKE Turns 10 hackathon](https://gketurns10.devpost.com/). We're creating a living virtual world where AI agents simulate real economic behaviors - shopping, working, and managing businesses - all interacting with production microservices (Online Boutique and Bank of Anthos) deployed on Google Kubernetes Engine.

### The Vision

Imagine a virtual city where:
- **AI Customers** wake up, check their bank balance, and shop for necessities
- **AI Employees** work shifts at the boutique or bank, earning salaries
- **AI Managers** make inventory decisions and set pricing strategies
- **Time flows 60x faster** - every real minute is an hour in their world

All powered by the Agent Development Kit (ADK) and communicating via the A2A protocol, demonstrating how AI can enhance existing microservices without modifying their code.

## 🏗️ Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                     AI Agent Ecosystem                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Customer │  │ Employee │  │ Manager  │  │   More   │   │
│  │  Agents  │  │  Agents  │  │  Agents  │  │  Agents  │   │
│  └─────┬────┘  └─────┬────┘  └─────┬────┘  └─────┬────┘   │
│        └──────────────┴──────────────┴──────────────┘       │
│                            │                                 │
│                    A2A Protocol Layer                        │
└────────────────────────────┬─────────────────────────────────┘
                             │
┌────────────────────────────┼─────────────────────────────────┐
│                     Service Layer                             │
│                            │                                  │
│  ┌─────────────────────────┴──────────────────────────┐     │
│  │            Payment Integration Service              │     │
│  │         (Bridges Boutique ←→ Bank Transfers)       │     │
│  └─────────────┬──────────────────┬───────────────────┘     │
│                │                  │                          │
│  ┌─────────────▼──────┐  ┌───────▼──────────┐              │
│  │  Online Boutique   │  │  Bank of Anthos  │              │
│  │   (E-commerce)     │  │    (Banking)     │              │
│  └────────────────────┘  └──────────────────┘              │
│                                                              │
│                    Google Kubernetes Engine                  │
└──────────────────────────────────────────────────────────────┘
```

### Service Integration

1. **Online Boutique**: E-commerce platform where agents purchase food, clothing, and entertainment
2. **Bank of Anthos**: Banking system managing agent accounts, salaries, and transactions
3. **Payment Integration**: Custom gRPC bridge enabling real bank transfers for purchases

### Agent Hierarchy

| Agent Type | Access Level | Capabilities |
|------------|--------------|--------------|
| **Customers** | Frontend only | Browse products, make purchases, check balance |
| **Employees** | Backend services | Process orders, handle transactions, customer service |
| **Managers** | Administrative | Set prices, manage inventory, business analytics |

## 🛠️ Technical Stack

- **Platform**: Google Kubernetes Engine (GKE)
- **AI Framework**: Agent Development Kit (ADK)
- **Communication**: A2A Protocol (https://a2a-protocol.org/dev/)
- **GitOps**: Argo CD + Kustomize
- **Languages**: Go (payment service), Python/Java (agents - planned)
- **Monitoring**: Prometheus metrics, structured logging


# Local development
# Configure Tilt
cp tilt_config.json.example tilt_config.json
# Edit tilt_config.json with your settings

# Start development environment
tilt up
```

## 🎮 Agent Simulation Concept

### Time Management
- **Virtual Time**: 1 real minute = 1 virtual hour
- **Work Schedule**: Agents work 8-hour shifts (8 real minutes)
- **Shopping Patterns**: Peak hours, seasonal variations

### Economic Behaviors
- **Income**: Agents earn salaries based on roles
- **Spending**: Necessity purchases (food) vs discretionary (entertainment)
- **Savings**: Agents maintain emergency funds
- **Credit**: Can take loans for large purchases

### Agent Tools
Each agent type gets specific tools:
- **Customers**: `browse_products()`, `add_to_cart()`, `checkout()`, `check_balance()`
- **Employees**: `process_order()`, `update_inventory()`, `assist_customer()`
- **Managers**: `set_price()`, `order_stock()`, `view_analytics()`, `hire_employee()`

## 📄 License

UNLICENSE

