# 7. Non-Functional Requirements (NFR) — Business-Grade Detail

NFRs are non-negotiable business constraints because they define customer experience, reliability of money movement, regulatory readiness, and operational cost.

## 7.1 NFR Catalog (Derived from SRS Section 5)

### NFR-046 — Database Technology

- **Business requirement:** Maintain relational data integrity using **PostgreSQL V17** with JSONB support
- **Target/Measurement:** ACID compliance tests
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-047 — Caching & Session

- **Business requirement:** Use **Redis** for distributed caching and session management
- **Target/Measurement:** Cache hit ratio monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-048 — API Protocol

- **Business requirement:** Microservices communication shall use **gRPC with Protocol Buffers** (Category 1)
- **Target/Measurement:** Inter-service latency metrics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-049 — Client API

- **Business requirement:** Client-facing APIs shall use **REST (OpenAPI 3.0)** with **JWT** authentication
- **Target/Measurement:** Schema validation, Token checks
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-050 — Public Integration

- **Business requirement:** External integrations shall use **RESTful APIs** with **OpenAPI 3.0** specifications (Category 3)
- **Target/Measurement:** Swagger validator pass
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-051 — Search Engine

- **Business requirement:** Full-text search capabilities shall be implemented using **PostgreSQL Full-Text Search** or dedicated engine
- **Target/Measurement:** Query performance <200ms
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-052 — Object Storage

- **Business requirement:** Document and static asset storage shall use **S3-compatible storage** (AWS/DigitalOcean)
- **Target/Measurement:** Upload/Download latency
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-053 — Message Broker

- **Business requirement:** Asynchronous event processing shall be handled by **Apache Kafka**
- **Target/Measurement:** Throughput monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-054 — Time-Series Data

- **Business requirement:** IoT telemetry data shall be stored in **TimescaleDB**
- **Target/Measurement:** Ingestion rate monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-055 — Vector Database

- **Business requirement:** Vector embeddings for AI features shall be stored in **Pgvector** or **Pinecone**
- **Target/Measurement:** Similarity search latency
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-056 — Graph Database

- **Business requirement:** Fraud visualization and relationship mapping shall use **Neo4j** or **Amazon Neptune**
- **Target/Measurement:** Graph traversal depth/speed
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-057 — Columnar Database

- **Business requirement:** High-performance analytics shall use **ClickHouse** or **Druid**
- **Target/Measurement:** Analytical query speed
- **Priority:** D
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-058 — Financial Ledger

- **Business requirement:** Double-entry bookkeeping shall be enforced using **TigerBeetle**
- **Target/Measurement:** Ledger reconciliation check
- **Priority:** M3
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-059 — Mobile Framework

- **Business requirement:** Cross-platform mobile application shall be built using **React Native**
- **Target/Measurement:** Code reuse >80%
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-060 — CDN & Security

- **Business requirement:** Public entry points shall be secured via **Cloudflare** proxy
- **Target/Measurement:** WAF block rate
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-001 — API response time for policy operations

- **Business requirement:** < 500ms (95th percentile)
- **Target/Measurement:** Application performance monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-002 — Database query response time

- **Business requirement:** < 100ms (average)
- **Target/Measurement:** Database monitoring tools
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-003 — Mobile app startup time

- **Business requirement:** < 3 seconds
- **Target/Measurement:** App performance analytics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-004 — Web portal page load time

- **Business requirement:** < 2 seconds
- **Target/Measurement:** Browser performance tools
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-005 — Payment processing time

- **Business requirement:** < 10 seconds end-to-end
- **Target/Measurement:** Payment gateway analytics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-006 — Claim processing automation

- **Business requirement:** 80% straight-through processing
- **Target/Measurement:** Business process monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-007 — Report generation time

- **Business requirement:** < 30 seconds for standard reports
- **Target/Measurement:** Reporting system metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-008 — Search functionality response

- **Business requirement:** < 200ms for basic searches
- **Target/Measurement:** Search performance monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-009 — Concurrent user support

- **Business requirement:** 10,000 active users
- **Target/Measurement:** Load testing and monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-010 — Transaction throughput

- **Business requirement:** 1,000 TPS (policies + claims)
- **Target/Measurement:** Performance testing
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-011 — Database scalability

- **Business requirement:** 100 million policy records
- **Target/Measurement:** Database performance testing
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-012 — Auto-scaling capability

- **Business requirement:** Scale out/in based on load
- **Target/Measurement:** Infrastructure monitoring
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-013 — Peak load handling

- **Business requirement:** 5x normal load during campaigns
- **Target/Measurement:** Stress testing
- **Priority:** M3
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-014 — Storage scalability

- **Business requirement:** 10TB+ document storage
- **Target/Measurement:** Cloud storage metrics
- **Priority:** M3
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-015 — System availability

- **Business requirement:** 99.5% uptime (M1), 99.9% (M2)
- **Target/Measurement:** Infrastructure monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-016 — Recovery Time Objective (RTO)

- **Business requirement:** 4 hours maximum
- **Target/Measurement:** Disaster recovery testing
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-017 — Recovery Point Objective (RPO)

- **Business requirement:** 1 hour maximum data loss
- **Target/Measurement:** Backup and recovery testing
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-018 — Mean Time To Recovery (MTTR)

- **Business requirement:** < 2 hours
- **Target/Measurement:** Incident response metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-019 — Service degradation handling

- **Business requirement:** Graceful degradation during outages
- **Target/Measurement:** Chaos engineering testing
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-020 — Data backup frequency

- **Business requirement:** Real-time replication + daily backups
- **Target/Measurement:** Backup monitoring
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-029 — User satisfaction score

- **Business requirement:** 4.5+ stars on app stores
- **Target/Measurement:** User feedback and ratings
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-030 — Task completion rate

- **Business requirement:** 95% for critical user journeys
- **Target/Measurement:** User experience analytics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-031 — Learning curve

- **Business requirement:** New users complete first task < 5 minutes
- **Target/Measurement:** User onboarding metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-032 — Error recovery

- **Business requirement:** Clear error messages with action guidance
- **Target/Measurement:** Error tracking and analysis
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-033 — Accessibility compliance

- **Business requirement:** WCAG 2.1 AA compliance
- **Target/Measurement:** Accessibility testing tools
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-034 — Multi-language support

- **Business requirement:** Bengali and English localization
- **Target/Measurement:** Localization testing
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-035 — Code coverage

- **Business requirement:** 80% unit test coverage
- **Target/Measurement:** Automated testing reports
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-036 — Deployment frequency

- **Business requirement:** Daily deployments capability
- **Target/Measurement:** CI/CD pipeline metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-037 — Mean Time To Deploy

- **Business requirement:** < 30 minutes for hotfixes
- **Target/Measurement:** Deployment automation metrics
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-038 — Monitoring coverage

- **Business requirement:** 100% critical path monitoring
- **Target/Measurement:** Observability platform
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-039 — Log aggregation

- **Business requirement:** Centralized logging for all services
- **Target/Measurement:** Logging platform metrics
- **Priority:** M1
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

### NFR-040 — Documentation currency

- **Business requirement:** API documentation auto-generated
- **Target/Measurement:** Documentation automation
- **Priority:** M2
- **Why this matters (business):** prevents customer drop-off, failed payments, slow claims, and audit risk.

[[[PAGEBREAK]]]
