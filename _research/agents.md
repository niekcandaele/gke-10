# AI Employee Agents for GKE Hackathon

## 🔌 A2A Protocol Integration

The Agent-to-Agent (A2A) protocol enables our agents to collaborate across platform boundaries without direct API access. Instead of boutique agents calling bank APIs (impossible due to security), agents request services from each other through standardized A2A messages.

### How It Works
- Each agent advertises capabilities via an **Agent Card**
- Agents discover what other agents can do
- Agents request services through A2A messages
- Responses provide insights, not raw data
- All communications are logged for audit

### Key Benefits
- **Security**: Bank data stays with bank agents
- **Privacy**: Agents share insights, not customer data  
- **Flexibility**: New agents integrate without API changes
- **Natural**: Agents communicate in task-oriented messages

## 🏦 Bank of Anthos Employee Agents

### 1. Fraud Detection Analyst
**Role:** Monitor and prevent fraudulent transactions across banking and e-commerce platforms

**Agent Card:**
```json
{
  "agent": "fraud-detection-analyst",
  "platform": "bank-of-anthos",
  "capabilities": [
    "check-fraud-risk",
    "freeze-account",
    "analyze-transaction-pattern",
    "validate-purchase"
  ],
  "accepts": ["user-id", "transaction-id", "amount", "pattern-type"],
  "provides": ["risk-score", "fraud-alert", "account-status"]
}
```

**Direct API Calls:**
- Calls: `Ledger Writer Service` - Read incoming transactions
- Calls: `Transaction History Service` - Analyze historical patterns
- Calls: `Balance Reader Service` - Check for unusual balance changes
- Calls: `User Service` - Get user profile and update account status
- Calls: `Contacts Service` - Verify trusted contacts

**A2A Interactions:**
- **Receives from Personal Shopping Assistant**: "Is this purchase pattern normal for user X?"
- **Receives from Payment Reconciliation**: "Investigate failed payment for fraud indicators"
- **Broadcasts to all agents**: "ALERT: Freeze all activity for user X - suspected compromise"
- **Responds to Customer Support**: "Account frozen due to: [specific pattern detected]"

**Key Responsibilities:**
- Monitor transaction patterns and velocity
- Cross-reference shopping behavior with banking patterns
- Coordinate platform-wide security responses
- Generate risk scores for other agents

### 2. Customer Service Representative
**Role:** Provide automated customer support for banking and payment processing issues

**Responsibilities & API Calls:**
- **Handle login issues and password resets**
  - Calls: `User Service` - Reset passwords, unlock accounts
  - Calls: `Frontend Service` - Generate reset tokens
- **Answer account balance inquiries**
  - Calls: `Balance Reader Service` - Get current balance
  - Calls: `Transaction History Service` - Recent transaction summary
  - Shows pending boutique purchases affecting available balance
- **Process dispute claims**
  - Calls: `Transaction History Service` - Retrieve transaction details
  - Calls: `Ledger Writer Service` - Submit reversal transactions
  - Creates: Dispute tickets in custom database
  - Handles boutique purchase disputes and payment failures
  - Coordinates with Payment Clearing House for transaction reversals
- **Manage contact lists**
  - Calls: `Contacts Service` - Add/remove/update contacts
  - Calls: `User Service` - Verify contact ownership

### 3. Financial Advisor
**Role:** Provide personalized financial insights including e-commerce spending analysis

**Agent Card:**
```json
{
  "agent": "financial-advisor",
  "platform": "bank-of-anthos",
  "capabilities": [
    "check-spending-power",
    "analyze-spending-patterns",
    "calculate-budget",
    "assess-affordability"
  ],
  "accepts": ["user-id", "amount", "purchase-type", "time-period"],
  "provides": ["available-budget", "spending-analysis", "affordability-score"]
}
```

**Direct API Calls:**
- Calls: `Transaction History Service` - Get categorized transactions
- Calls: `Balance Reader Service` - Track balance trends and savings
- Calls: `User Service` - Get user preferences

**A2A Interactions:**
- **Receives from Personal Shopping Assistant**: "Can user X afford $500 purchase?"
- **Responds with**: "User has $800 discretionary budget, purchase is 62.5% of available funds"
- **Receives from Marketing Specialist**: "Which users have >$1000 available for shopping?"
- **Notifies Marketing Specialist**: "User X received salary, discretionary budget increased"
- **Alerts Customer Support**: "User Y has insufficient funds for pending purchase"

**Key Responsibilities:**
- Calculate real-time spending power for shopping agents
- Identify optimal shopping times based on cash flow
- Provide affordability assessments without exposing exact balances
- Alert agents when spending patterns indicate issues

### 4. Compliance Officer
**Role:** Ensure regulatory compliance across banking and e-commerce transactions

**Responsibilities & API Calls:**
- **Monitor AML patterns**
  - Calls: `Ledger Database` (read-only) - Complex transaction queries
  - Calls: `Transaction History Service` - Pattern matching
  - Calls: `User Service` - User risk profiles
  - Monitors boutique purchases for money laundering patterns
  - Detects rapid buy-return cycles that could indicate fraud
- **Generate regulatory reports**
  - Calls: `Ledger Database` (read-only) - Aggregate reporting data
  - Calls: `Balance Reader Service` - Account summaries
  - Creates: Compliance reports in custom database
  - Includes e-commerce transaction volumes in regulatory filings
- **Ensure KYC requirements**
  - Calls: `User Service` - Verify user documentation
  - Creates: KYC verification records
  - Validates boutique user identity matches bank account holder
- **Flag threshold violations**
  - Calls: `Ledger Writer Service` - Monitor large transactions
  - Calls: `User Service` - Update risk flags
  - Flags high-value boutique purchases requiring additional verification

### 5. Account Manager
**Role:** Handle account lifecycle management and customer relationships

**Responsibilities & API Calls:**
- **Onboard new customers**
  - Calls: `User Service` - Create new accounts
  - Calls: `Accounts Database` - Initialize account data
  - Calls: `Contacts Service` - Set up initial contacts
- **Manage account upgrades**
  - Calls: `User Service` - Update account tiers
  - Calls: `Balance Reader Service` - Verify eligibility
- **Handle account closures**
  - Calls: `Balance Reader Service` - Check final balance
  - Calls: `User Service` - Deactivate account
  - Calls: `Transaction History Service` - Generate final statement
- **Coordinate beneficiaries**
  - Calls: `Contacts Service` - Manage beneficiary lists
  - Calls: `User Service` - Update beneficiary permissions

---

## 🛍️ Online Boutique Employee Agents

### 1. Personal Shopping Assistant
**Role:** Provide personalized shopping experiences with real-time budget awareness

**Agent Card:**
```json
{
  "agent": "personal-shopping-assistant",
  "platform": "online-boutique",
  "capabilities": [
    "recommend-products",
    "check-affordability",
    "create-collections",
    "validate-purchase"
  ],
  "accepts": ["user-id", "budget", "preferences", "cart-contents"],
  "provides": ["product-recommendations", "budget-collections", "affordability-check"]
}
```

**Direct API Calls:**
- Calls: `RecommendationService.ListRecommendations()` - Base recommendations
- Calls: `ProductCatalogService.ListProducts()` - Full catalog access
- Calls: `CartService.GetCart()` - Current cart context
- Calls: `CurrencyService.Convert()` - Multi-currency support
- Calls: `ShippingService.GetQuote()` - Include shipping in budget

**A2A Interactions:**
- **Requests from Financial Advisor**: "What is user X's available shopping budget?"
- **Receives response**: "User has $500 discretionary funds available"
- **Requests from Fraud Analyst**: "Is this shopping pattern normal for user X?"
- **Receives from Marketing Specialist**: "Create payday collection for users Y, Z"
- **Notifies Customer Support**: "User attempted purchase exceeding available funds"

**Key Responsibilities:**
- Filter recommendations based on affordability (via A2A to Financial Advisor)
- Create "Within Your Budget" collections using bank insights
- Warn users before they exceed their budget
- Coordinate with bank agents for purchase validation

### 2. Customer Support Specialist
**Role:** Handle post-purchase support with emphasis on bank payment issues

**Responsibilities & API Calls:**
- **Track order status**
  - Calls: `CheckoutService.PlaceOrder()` results - Order details
  - Creates: Order tracking in custom database
  - Links order status to bank transaction ID
- **Handle shipping inquiries**
  - Calls: `ShippingService.GetQuote()` - Shipping estimates
  - Calls: `ShippingService.ShipOrder()` status - Tracking info
- **Process returns/refunds**
  - Calls: `Payment Clearing House` - Process bank account refunds
  - Calls: `EmailService.SendOrderConfirmation()` - Send return confirmations
  - Creates: Return records in custom database
  - Coordinates with Bank Customer Service for payment reversals
- **Cart abandonment recovery**
  - Calls: `CartService.GetCart()` - Retrieve abandoned carts
  - Calls: `EmailService.SendOrderConfirmation()` - Send recovery emails
  - Calls: `ProductCatalogService.GetProduct()` - Check item availability
  - Identifies if abandonment was due to insufficient funds

### 3. Inventory Manager
**Role:** Monitor and manage product availability

**Responsibilities & API Calls:**
- **Monitor product availability**
  - Calls: `ProductCatalogService.ListProducts()` - Current inventory
  - Creates: Inventory tracking in custom database
- **Low stock alerts**
  - Calls: `ProductCatalogService.GetProduct()` - Stock levels
  - Creates: Alert notifications for admin
- **Suggest alternatives**
  - Calls: `ProductCatalogService.SearchProducts()` - Find similar items
  - Calls: `RecommendationService.ListRecommendations()` - Related products
- **Restocking recommendations**
  - Analyzes: Cart and checkout patterns
  - Creates: Restocking reports in custom database

### 4. Marketing Specialist
**Role:** Create targeted campaigns based on actual customer spending power

**Agent Card:**
```json
{
  "agent": "marketing-specialist",
  "platform": "online-boutique",
  "capabilities": [
    "create-promotions",
    "analyze-abandonment",
    "segment-customers",
    "time-campaigns"
  ],
  "accepts": ["user-segment", "cart-data", "budget-range"],
  "provides": ["targeted-promotions", "recovery-campaigns", "customer-segments"]
}
```

**Direct API Calls:**
- Calls: `AdService.GetAds()` - Coordinate with ad system
- Calls: `CartService.GetCart()` - Cart analysis and abandonment
- Calls: `EmailService.SendOrderConfirmation()` - Campaign emails
- Calls: `RecommendationService.ListRecommendations()` - Cross-sell opportunities

**A2A Interactions:**
- **Requests from Financial Advisor**: "Which users received payday deposits today?"
- **Receives list**: "Users A, B, C have increased discretionary budgets"
- **Requests from Financial Advisor**: "Why might user X have abandoned cart?"
- **Receives**: "Insufficient funds - cart total exceeds available budget by $50"
- **Notifies Personal Shopping Assistant**: "Launch payday promotions for user segment"
- **Coordinates with Account Manager (Bank)**: "Identify VIP customers for premium campaign"

**Key Responsibilities:**
- Time campaigns based on payday cycles (via A2A insights)
- Create "You can afford this!" campaigns using budget data
- Identify cart abandonment due to insufficient funds
- Segment customers by actual spending power, not just behavior

### 5. Pricing Analyst
**Role:** Optimize pricing strategies and manage discounts

**Responsibilities & API Calls:**
- **Dynamic pricing adjustments**
  - Calls: `ProductCatalogService.GetProduct()` - Current prices
  - Calls: `CurrencyService.GetSupportedCurrencies()` - Regional pricing
  - Creates: Price adjustment rules in custom database
- **Apply personalized discounts**
  - Calls: `CheckoutService.PlaceOrder()` - Apply at checkout
  - Calls: `CartService.GetCart()` - Calculate cart discounts
- **Monitor currency fluctuations**
  - Calls: `CurrencyService.Convert()` - Exchange rate tracking
  - Calls: `CurrencyService.GetSupportedCurrencies()` - Available currencies
- **Optimize shipping costs**
  - Calls: `ShippingService.GetQuote()` - Analyze shipping patterns
  - Creates: Shipping optimization rules

### 6. Quality Assurance Specialist
**Role:** Monitor product quality and customer satisfaction

**Responsibilities & API Calls:**
- **Monitor customer reviews**
  - Creates: Review system in custom database
  - Calls: `ProductCatalogService.GetProduct()` - Link reviews to products
- **Track return rates**
  - Analyzes: Return records from Customer Support agent
  - Calls: `ProductCatalogService.GetProduct()` - Product return statistics
- **Identify quality issues**
  - Aggregates: Customer feedback patterns
  - Calls: `EmailService.SendOrderConfirmation()` - Quality surveys
- **Suggest improvements**
  - Calls: `ProductCatalogService.ListProducts()` - Product analysis
  - Creates: Quality reports in custom database

---

## 👥 Unified Customer Persona Agents

Each persona represents a customer who actively uses BOTH Bank of Anthos AND Online Boutique, with interconnected financial and shopping behaviors.

### 1. Priya - The Strategic Saver & Planner
**Profile:** Software engineer in Bangalore who optimizes both finances and purchases

**Banking Behavior (Bank of Anthos):**
- Deposits ₹150,000 salary on 1st of each month
- Auto-transfers 40% to savings/investments
- Weekly balance checks (Sunday evenings IST)
- Systematic investment plans (SIPs)
- Sets aside 10% for planned purchases

**Shopping Behavior (Online Boutique):**
- Shops only during planned purchase windows
- Researches products for weeks before buying
- Waits for Diwali/year-end sales
- Bulk buys essentials quarterly
- Never impulse purchases

**Interconnected Patterns:**
- Shopping budget auto-transfers to separate account after salary
- Cart abandonment if item exceeds allocated budget
- Increased shopping only after bonuses/increments
- Reviews Financial Advisor recommendations before major purchases
- Shopping receipts categorized for tax planning

**API Interactions:**
- Bank: `User Service`, `Ledger Writer Service`, `Balance Reader Service`
- Boutique: `ProductCatalogService.SearchProducts()`, `CartService.AddItem()`, `CheckoutService.PlaceOrder()`
- Cross-platform: Checks balance before checkout completion

**Employee Agent Touchpoints:**
- Financial Advisor (Bank) influences Personal Shopping Assistant (Boutique) recommendations
- Marketing Specialist (Boutique) promotions align with her savings goals
- Quality Assurance (Boutique) receives her detailed product reviews

### 2. João - The Global Entrepreneur
**Profile:** Import/export business owner in São Paulo managing B2B and personal accounts

**Banking Behavior (Bank of Anthos):**
- Multiple daily international transfers (BRL/USD/EUR/CNY)
- Business account with R$500K+ monthly flow
- Personal account for family expenses
- Friday bulk payments to suppliers
- Currency hedging operations

**Shopping Behavior (Online Boutique):**
- Sources product samples for import opportunities
- Bulk purchases for business gifts
- Personal luxury purchases after deals close
- Ships to international business partners
- Tests international shipping logistics

**Interconnected Patterns:**
- Business deal closures trigger luxury shopping
- Uses boutique for market research (trending products)
- Currency fluctuations affect shopping decisions
- Business expense purchases through boutique
- Cash flow determines shopping volume

**API Interactions:**
- Bank: `Ledger Writer Service` (high-frequency), `User Service` (multi-session), `Balance Reader Service`
- Boutique: `CurrencyService.Convert()`, `ShippingService.GetQuote()`, `CheckoutService.PlaceOrder()` (bulk)
- Cross-platform: Business expenses tracked across both

**Employee Agent Touchpoints:**
- Compliance Officer (Bank) ↔ Inventory Manager (Boutique) for bulk orders
- Account Manager (Bank) coordinates with Personal Shopping Assistant (Boutique) for VIP treatment
- Currency movements trigger Pricing Analyst (Boutique) alerts

### 3. Amara - The Budget-Conscious Student
**Profile:** Masters student in London from Nigeria, managing tight finances

**Banking Behavior (Bank of Anthos):**
- £800-1200 monthly part-time income (irregular)
- £200 monthly remittance to Lagos
- £50 weekly grocery budget
- Overdraft protection active
- Student loan payments

**Shopping Behavior (Online Boutique):**
- Only shops with student discounts
- Abandons cart if total exceeds weekly budget
- Compares prices for hours
- Ships gifts home during holidays
- Returns items if unexpected expenses arise

**Interconnected Patterns:**
- Checks bank balance before every purchase
- Shopping restricted to post-payday week
- Returns increase near rent due date
- Remittance amount affects shopping budget
- Financial stress visible in both platforms

**API Interactions:**
- Bank: `Balance Reader Service` (constant), `Ledger Writer Service` (remittances)
- Boutique: `CartService.GetCart()` (abandoned carts), `PaymentService.Charge()` (refunds frequent)
- Cross-platform: Real-time balance checks during checkout

**Employee Agent Touchpoints:**
- Financial Advisor (Bank) helps optimize for shopping
- Customer Support (Boutique) handles frequent returns
- Marketing Specialist (Boutique) sends student discount alerts

### 4. Wei - The Restaurant Owner
**Profile:** Vancouver Chinatown restaurant owner managing business and family

**Banking Behavior (Bank of Anthos):**
- Daily cash deposits at 11 PM PST (CAD 3-5K)
- Weekly supplier payments (CAD 15-20K)
- Bi-weekly payroll for 8 employees
- Business and personal accounts
- Seasonal fluctuations (Chinese New Year 3x normal)

**Shopping Behavior (Online Boutique):**
- Purchases restaurant uniforms/merchandise
- Bulk buys during supplier sales
- Personal shopping for family of 4
- Ships gifts to family in Guangzhou
- Holiday decorations for restaurant

**Interconnected Patterns:**
- Good revenue weeks trigger family shopping
- Supplier payment timing affects boutique purchases
- Uses boutique for restaurant promotional items
- Cash flow determines inventory stocking
- Business seasonality drives both platforms

**API Interactions:**
- Bank: `Ledger Writer Service` (daily deposits), `User Service` (business features)
- Boutique: `CheckoutService.PlaceOrder()` (bulk), `ShippingService.ShipOrder()` (international)
- Cross-platform: Business expense tracking

**Employee Agent Touchpoints:**
- Compliance Officer (Bank) monitors cash deposits
- Inventory Manager (Boutique) coordinates bulk orders
- Account Manager (Bank) and Personal Shopping Assistant (Boutique) provide business support

### 5. Carlos - The Gig Economy Worker
**Profile:** Mexico City rideshare/delivery driver with variable income

**Banking Behavior (Bank of Anthos):**
- Daily micro-deposits (MXN 200-1500)
- Instant transfers for gas/expenses
- Weekly earnings vary 50%+
- No savings buffer
- Frequent overdrafts

**Shopping Behavior (Online Boutique):**
- Only shops on high-earning days
- Abandons cart if daily earnings low
- Purchases work essentials (phone accessories)
- Birthday/holiday gifts when possible
- Returns items during slow weeks

**Interconnected Patterns:**
- Shopping directly correlated to daily earnings
- Friday night earnings trigger weekend shopping
- Cart abandonment rate 70%+
- Uses boutique purchases as earning motivation
- Financial volatility visible across platforms

**API Interactions:**
- Bank: `Ledger Writer Service` (micro-transactions), `Balance Reader Service` (constant)
- Boutique: `CartService.AddItem()` then `EmptyCart()` frequently
- Cross-platform: Real-time earning-to-shopping pipeline

**Employee Agent Touchpoints:**
- Customer Service (Bank) handles overdraft issues
- Marketing Specialist (Boutique) times promotions to payday
- Financial Advisor (Bank) struggles with irregular income patterns

### 6. Fatima - The Family Financial Manager
**Profile:** Dubai teacher managing household finances and shopping

**Banking Behavior (Bank of Anthos):**
- AED 15,000 monthly salary
- Separate accounts for household, savings, family support
- Sends money to parents in Jordan
- Children's education savings plan
- Tracks every expense meticulously

**Shopping Behavior (Online Boutique):**
- Bulk family shopping during sales
- Coordinates Ramadan/Eid gifts for extended family
- School supplies in August
- Compares prices across multiple sessions
- Ships to family across Middle East

**Interconnected Patterns:**
- Shopping budget predetermined by banking allocations
- Eid bonuses trigger family gift shopping
- Education expenses affect shopping patterns
- Savings goals limit discretionary shopping
- Family remittances reduce shopping budget

**API Interactions:**
- Bank: `User Service` (multiple accounts), `Transaction History Service` (expense tracking)
- Boutique: `ShippingService.GetQuote()` (multiple addresses), `CurrencyService.Convert()` (AED/USD/JOD)
- Cross-platform: Budget enforcement across platforms

**Employee Agent Touchpoints:**
- Financial Advisor (Bank) helps balance family obligations
- Personal Shopping Assistant (Boutique) suggests family deals
- Marketing Specialist (Boutique) targets cultural holidays

### 7. Yuki - The Design Professional
**Profile:** Tokyo architect balancing aesthetic preferences with financial discipline

**Banking Behavior (Bank of Anthos):**
- ¥600,000 monthly salary
- 30% to investment portfolio
- Separate account for design purchases
- Annual bonuses (2x salary)
- International client payments

**Shopping Behavior (Online Boutique):**
- Researches products for weeks
- Purchases high-quality, sustainable items
- Willing to pay premium for craftsmanship
- Seasonal wardrobe updates
- Gifts for clients

**Interconnected Patterns:**
- Bonus payments trigger luxury purchases
- Investment performance affects shopping mood
- Business expenses through boutique
- Quality over quantity aligns with savings goals
- International projects drive currency/shipping needs

**API Interactions:**
- Bank: `Balance Reader Service`, `Transaction History Service` (categorization)
- Boutique: `ProductCatalogService.GetProduct()` (detailed inspection), `RecommendationService`
- Cross-platform: Business expense reconciliation

**Employee Agent Touchpoints:**
- Financial Advisor (Bank) supports investment strategy
- Quality Assurance (Boutique) receives detailed feedback
- Personal Shopping Assistant (Boutique) provides curation

### 8. Dmitri - The Minimalist Engineer
**Profile:** Moscow software engineer focused on functionality

**Banking Behavior (Bank of Anthos):**
- ₽200,000 monthly salary
- 60% automated to savings/investments
- Minimal transaction count
- Annual travel fund
- Cryptocurrency investments

**Shopping Behavior (Online Boutique):**
- Buys only when items break/wear out
- Focuses on durability and specifications
- Bulk purchases to minimize shopping frequency
- Standard shipping always
- No emotional purchases

**Interconnected Patterns:**
- Shopping inversely correlated with investment performance
- Annual shopping budget pre-allocated
- Purchase timing based on bank balance thresholds
- Returns rare due to extensive research
- Efficiency focus across both platforms

**API Interactions:**
- Bank: `User Service` (minimal sessions), `Ledger Writer Service` (automated)
- Boutique: `ProductCatalogService.SearchProducts()` (specific searches), `CheckoutService.PlaceOrder()` (quick)
- Cross-platform: Minimal touchpoints

**Employee Agent Touchpoints:**
- Ignores most employee agents
- Only contacts Customer Support for defects
- Financial Advisor (Bank) automation recommendations adopted

### 9. Maria - The Community Organizer
**Profile:** Mexico City event planner managing personal and community finances

**Banking Behavior (Bank of Anthos):**
- MXN 25,000 monthly income
- Manages community fund account
- Coordinates group payments
- Event deposit management
- Transparent transaction records

**Shopping Behavior (Online Boutique):**
- Bulk purchases for events (quinceañeras, weddings)
- Negotiates group discounts
- Ships to multiple addresses
- Time-sensitive deliveries
- Returns/exchanges for size issues

**Interconnected Patterns:**
- Event deposits trigger shopping sprees
- Community funds used for group purchases
- Personal shopping limited by event schedule
- Banking transparency required for community trust
- Seasonal patterns (wedding season, holidays)

**API Interactions:**
- Bank: `Contacts Service` (large network), `Transaction History Service` (reporting)
- Boutique: `CheckoutService.PlaceOrder()` (bulk with multiple addresses), `EmailService` (coordination)
- Cross-platform: Fund collection to purchase pipeline

**Employee Agent Touchpoints:**
- Account Manager (Bank) helps with community features
- Customer Support (Boutique) manages complex deliveries
- Marketing Specialist (Boutique) provides group discounts

### 10. Kwame - The Tech Innovator
**Profile:** Lagos startup founder testing platforms while building his business

**Banking Behavior (Bank of Anthos):**
- Variable income (₦0-5M monthly)
- Multiple currency accounts (NGN/USD/EUR)
- Angel investment deposits
- High-risk investment profile
- International wire transfers

**Shopping Behavior (Online Boutique):**
- Beta tests every feature
- Studies recommendation algorithms
- Mobile-first power user
- Creates video reviews
- Reports bugs enthusiastically

**Interconnected Patterns:**
- Funding rounds trigger tech equipment shopping
- Tests payment methods across platforms
- Uses both platforms for startup research
- Financial volatility drives shopping patterns
- Innovation focus across both services

**API Interactions:**
- Bank: All services (edge case testing)
- Boutique: All services (comprehensive testing)
- Cross-platform: Studies integration possibilities

**Employee Agent Touchpoints:**
- Provides feedback to all employee agents
- Direct channel to Quality Assurance teams
- Tests agent response limits

### 🔄 Unified Behavioral Patterns

#### Payday Effects
- **Priya:** Triggers automatic savings and shopping budget allocation
- **Carlos:** Immediate shopping based on daily earnings
- **Amara:** Week-long shopping window opens
- **Wei:** Supplier payments before personal shopping
- **Fatima:** Family allocations before discretionary spending

#### Financial Stress Responses
- **Amara:** Cart abandonment and return rates increase
- **Carlos:** Shopping completely stops
- **Wei:** Shifts from personal to business-only purchases
- **Dmitri:** Unaffected due to pre-allocated budgets
- **Maria:** Prioritizes community over personal purchases

#### Seasonal Patterns
- **Fatima:** Ramadan/Eid banking and shopping surge
- **Wei:** Chinese New Year cash flow and gift shopping
- **Maria:** Wedding season fund collection and bulk purchasing
- **Priya:** Diwali bonuses and annual shopping
- **Yuki:** Year-end bonuses drive luxury purchases

#### Cross-Platform Triggers
- Low bank balance → Abandoned boutique carts
- Successful investment → Luxury shopping
- Business revenue → B2B boutique purchases
- Overdraft warning → Shopping freeze
- Bonus deposit → Wishlist purchases
- International transfer → Currency service usage
- Fraud alert → All platform activity pauses

#### Inter-Customer Dynamics
- **Maria** coordinates group purchases with community members
- **Amara** splits expenses with other students
- **João** sends business gifts to partners
- **Wei** employs **Carlos** for delivery services
- **Fatima** shares deals with family network

---

## 📨 A2A Workflow Examples

### 1. Real-Time Affordability Check
```
User browses $500 item in boutique
↓
Personal Shopping Assistant → (A2A) → Financial Advisor
  Request: {
    "task": "check-affordability",
    "user_id": "priya-strategic",
    "amount": 500,
    "type": "discretionary"
  }
↓
Financial Advisor checks bank balance, calculates budget
↓
Financial Advisor → (A2A) → Personal Shopping Assistant
  Response: {
    "affordable": true,
    "budget_impact": "62.5% of discretionary funds",
    "recommendation": "affordable but significant"
  }
↓
Personal Shopping Assistant shows: "✓ Within budget (62.5% of available funds)"
```

### 2. Payday Promotion Campaign
```
Bank detects salary deposits (Friday 9am)
↓
Account Manager → (A2A broadcast) → Marketing Specialist
  "Payday detected for users: [priya, joão, wei]"
↓
Marketing Specialist → (A2A) → Financial Advisor
  "Calculate shopping budgets for payday users"
↓
Financial Advisor → (A2A) → Marketing Specialist
  "Priya: $800, João: $2000, Wei: $1500 available"
↓
Marketing Specialist → (A2A) → Personal Shopping Assistant
  "Create personalized collections for payday users with budgets"
↓
Personal Shopping Assistant creates "Payday Treats Within Your Budget"
```

### 3. Cart Abandonment Investigation
```
User abandons $300 cart
↓
Marketing Specialist → (A2A) → Financial Advisor
  "Why did user carlos-gig abandon $300 cart?"
↓
Financial Advisor → (A2A) → Marketing Specialist
  "Insufficient funds: balance $250, cart $300"
↓
Marketing Specialist → (A2A) → Customer Support (Bank)
  "Can we offer payment assistance for carlos-gig?"
↓
Customer Support → (A2A) → Marketing Specialist
  "User eligible for overdraft protection, sending assistance email"
```

### 4. Cross-Platform Fraud Response
```
Fraud Analyst detects: 5 purchases in 10 minutes
↓
Fraud Analyst → (A2A broadcast ALL)
  "SECURITY ALERT: Freeze user amara-student - unusual velocity"
↓
All agents receive and act:
- Personal Shopping Assistant: Blocks cart checkout
- Customer Support (Bank): Locks online banking
- Customer Support (Boutique): Flags account
- Payment Reconciliation: Blocks pending transactions
↓
Fraud Analyst → (A2A) → Customer Service Rep
  "Contact user via registered phone for verification"
```

### 5. Payment Failure Investigation
```
Payment fails at checkout
↓
Customer Support (Boutique) → (A2A) → Payment Reconciliation
  "Investigate failed payment: transaction_id_789"
↓
Payment Reconciliation → (A2A) → Fraud Analyst
  "Any blocks for user fatima-family?"
↓
Fraud Analyst → (A2A) → Payment Reconciliation
  "No fraud blocks active"
↓
Payment Reconciliation → (A2A) → Financial Advisor
  "Check account status for fatima-family"
↓
Financial Advisor → (A2A) → Payment Reconciliation
  "Daily transaction limit reached ($5000)"
↓
Payment Reconciliation → (A2A) → Customer Support (Boutique)
  "Payment failed: daily limit exceeded, reset at midnight"
```

---

## 🔄 New Cross-Platform Agents

### 1. Payment Reconciliation Agent
**Role:** Monitor and reconcile payments between Bank of Anthos and Online Boutique

**Responsibilities & API Calls:**
- **Monitor payment flows**
  - Calls: `Payment Clearing House Database` - Track all transactions
  - Calls: `Bank Transaction API` - Verify bank-side records
  - Calls: `Boutique Order API` - Verify boutique-side records
  - Creates: Discrepancy reports and alerts
- **Handle payment failures**
  - Detects: Stuck or failed transactions
  - Initiates: Automatic retry or reversal workflows
  - Alerts: Both Bank and Boutique support teams
- **Generate reconciliation reports**
  - Aggregates: Daily/hourly payment statistics
  - Identifies: Patterns in payment failures
  - Reports: Success rates, average processing time
- **Manage reversals**
  - Coordinates: Between bank and boutique for refunds
  - Ensures: Atomic transaction reversal
  - Tracks: Reversal completion status

### 2. Cross-Platform Fraud Analyst
**Role:** Unified fraud detection across banking and shopping behaviors

**Responsibilities & API Calls:**
- **Real-time transaction monitoring**
  - Calls: `Payment Clearing House` - All payment attempts
  - Calls: `Bank Fraud API` - Risk scoring
  - Calls: `Boutique Order History` - Purchase patterns
  - Creates: Unified risk profiles
- **Detect cross-platform fraud patterns**
  - Identifies: Account takeover spanning both platforms
  - Monitors: Velocity of bank-to-boutique transactions
  - Flags: Unusual shopping after large deposits
- **Coordinate platform-wide blocks**
  - Freezes: Bank account AND boutique access simultaneously
  - Notifies: All relevant employee agents
  - Creates: Incident reports for investigation
- **Machine learning model training**
  - Feeds: Combined transaction data to ML models
  - Improves: Fraud detection accuracy over time
  - Adapts: To new fraud patterns

---

## 🔗 Agent Communication Matrix

### Who Talks to Whom (A2A Connections)

| From ↓ / To → | Fraud Analyst | Financial Advisor | Customer Service (Bank) | Personal Shopping | Customer Support (Boutique) | Marketing Specialist | Payment Reconciliation |
|----------------|---------------|-------------------|------------------------|------------------|----------------------------|---------------------|----------------------|
| **Fraud Analyst** | - | Verify spending | Alert account issues | Block purchases | Alert payment blocks | - | Investigate fraud |
| **Financial Advisor** | Check risk | - | Suggest assistance | Provide budgets | Alert low funds | Send payday info | - |
| **Customer Service (Bank)** | Request freeze | Get analysis | - | - | Coordinate refunds | - | Payment status |
| **Personal Shopping** | Check patterns | Request budgets | - | - | Failed purchase info | Receive campaigns | - |
| **Customer Support (Boutique)** | Check blocks | Verify funds | Request reversal | Get cart info | - | - | Investigate failure |
| **Marketing Specialist** | - | Get payday users | - | Launch campaigns | - | - | - |
| **Payment Reconciliation** | Check fraud | Verify balance | Get account status | - | Send failure reason | - | - |

### Communication Patterns

**Request-Response**: One agent asks, another responds
- Personal Shopping ↔ Financial Advisor (budget checks)
- Marketing ↔ Financial Advisor (payday notifications)

**Broadcast**: One agent alerts all
- Fraud Analyst → ALL (security alerts)
- Payment Reconciliation → Support teams (mass failures)

**Chain**: Multi-hop investigation
- Customer Support → Payment Reconciliation → Fraud Analyst → Financial Advisor

---

## 🔧 Implementation Architecture

### Agent Communication Flow
```
User Request → AI Agent → Existing Microservice APIs → Response Processing → User Response
                   ↓
            Agent Database
            (Custom State)
```

### Inter-Agent Collaboration via A2A Protocol

**Within-Platform Collaboration:**
- **Bank**: Fraud Analyst → (A2A) → Customer Service: "Account compromised, initiate contact"
- **Bank**: Financial Advisor → (A2A) → Account Manager: "User eligible for premium tier"
- **Boutique**: Inventory Manager → (A2A) → Marketing: "Stock low, create urgency campaign"

**Cross-Platform A2A Workflows:**

**Budget-Aware Shopping:**
```
Personal Shopping Assistant → Financial Advisor: "Check budget for user X"
Financial Advisor → Personal Shopping Assistant: "$500 available"
Personal Shopping Assistant filters products within budget
```

**Payday Marketing:**
```
Account Manager → Marketing Specialist: "Payday for users [A,B,C]"
Marketing → Financial Advisor: "Get shopping budgets"
Financial Advisor → Marketing: "A:$800, B:$1200, C:$500"
Marketing → Personal Shopping: "Create budget collections"
```

**Security Response:**
```
Fraud Analyst → (A2A Broadcast): "FREEZE user Y - all platforms"
All agents acknowledge and act within their domain
```

**Payment Investigation:**
```
Customer Support → Payment Reconciliation: "Why did payment fail?"
Payment Reconciliation → Fraud Analyst: "Check for blocks"
Fraud Analyst → Payment Reconciliation: "No blocks"
Payment Reconciliation → Financial Advisor: "Check balance"
Financial Advisor → Payment Reconciliation: "Insufficient funds"
Payment Reconciliation → Customer Support: "NSF - needs $50 more"
```

### Technology Stack
- **Container Runtime:** Google Kubernetes Engine (GKE)
- **AI Model:** Google Gemini
- **Agent Communication:** A2A (Agent-to-Agent) protocol for all inter-agent messages
- **Service Communication:** gRPC for microservice calls, REST for Bank APIs
- **State Management:** Separate Redis/PostgreSQL for agent-specific data
- **Observability:** OpenTelemetry integration with existing services

### Key Design Principles
1. **No Core Modification:** Agents only consume existing APIs
2. **Stateless Operations:** Agents maintain their own state separately
3. **Fault Tolerance:** Graceful degradation if services unavailable
4. **Scalability:** Each agent can scale independently on GKE
5. **Security:** JWT tokens for service authentication, following existing patterns