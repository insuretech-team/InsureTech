# 5. Non-Functional Requirements & Technical Constraints

## 5.1 Technology Constraints

| NFR ID | Constraint Area | Requirement | Measurement | Priority |
|--------|-----------------|-------------|-------------|----------|
| NFR-046 | Database Technology | The system shall maintain relational data integrity using **PostgreSQL V17** with JSONB support | ACID compliance tests | M1 |
| NFR-047 | Caching & Session | The system shall use **Redis** for distributed caching and session management | Cache hit ratio monitoring | M1 |
| NFR-048 | API Protocol | Microservices communication shall use **gRPC with Protocol Buffers** (Category 1) | Inter-service latency metrics | M1 |
| NFR-049 | Client API | Client-facing APIs shall use **REST (OpenAPI 3.0)** with **JWT** authentication | Schema validation, Token checks | M1 |
| NFR-050 | Public Integration | External integrations shall use **RESTful APIs** with **OpenAPI 3.0** specifications (Category 3) | Swagger validator pass | D |
| NFR-051 | Search Engine | Full-text search capabilities shall be implemented using **PostgreSQL Full-Text Search** or dedicated engine | Query performance <200ms | M1 |
| NFR-052 | Object Storage | Document and static asset storage shall use **S3-compatible storage** (AWS/DigitalOcean) | Upload/Download latency | M1 |
| NFR-053 | Message Broker | Asynchronous event processing shall be handled by **Apache Kafka** | Throughput monitoring | M1 |
| NFR-054 | Time-Series Data | IoT telemetry data shall be stored in **TimescaleDB** | Ingestion rate monitoring | M2 |
| NFR-055 | Vector Database | Vector embeddings for AI features shall be stored in **Pgvector** or **Pinecone** | Similarity search latency | D |
| NFR-056 | Graph Database | Fraud visualization and relationship mapping shall use **Neo4j** or **Amazon Neptune** | Graph traversal depth/speed | D |
| NFR-057 | Columnar Database | High-performance analytics shall use **ClickHouse** or **Druid** | Analytical query speed | D |
| NFR-058 | Financial Ledger | Double-entry bookkeeping shall be enforced using **TigerBeetle** | Ledger reconciliation check | M3 |
| NFR-059 | Mobile Framework | Cross-platform mobile application shall be built using **React Native** | Code reuse >80% | M1 |
| NFR-060 | CDN & Security | Public entry points shall be secured via **Cloudflare** proxy | WAF block rate | M1 |



## 5.2 Performance Requirements

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-001 | API response time for policy operations | < 500ms (95th percentile) | Application performance monitoring | M1 |
| NFR-002 | Database query response time | < 100ms (average) | Database monitoring tools | M1 |
| NFR-003 | Mobile app startup time | < 3 seconds | App performance analytics | M1 |
| NFR-004 | Web portal page load time | < 2 seconds | Browser performance tools | M1 |
| NFR-005 | Payment processing time | < 10 seconds end-to-end | Payment gateway analytics | M1 |
| NFR-006 | Claim processing automation | 80% straight-through processing | Business process monitoring | M2 |
| NFR-007 | Report generation time | < 30 seconds for standard reports | Reporting system metrics | M2 |
| NFR-008 | Search functionality response | < 200ms for basic searches | Search performance monitoring | M2 |

## 5.3 Scalability Requirements

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-009 | Concurrent user support | 10,000 active users | Load testing and monitoring | M1 |
| NFR-010 | Transaction throughput | 1,000 TPS (policies + claims) | Performance testing | M2 |
| NFR-011 | Database scalability | 100 million policy records | Database performance testing | M2 |
| NFR-012 | Auto-scaling capability | Scale out/in based on load | Infrastructure monitoring | M2 |
| NFR-013 | Peak load handling | 5x normal load during campaigns | Stress testing | M3 |
| NFR-014 | Storage scalability | 10TB+ document storage | Cloud storage metrics | M3 |

## 5.4 Availability & Reliability

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-015 | System availability | 99.5% uptime (M1), 99.9% (M2) | Infrastructure monitoring | M1 |
| NFR-016 | Recovery Time Objective (RTO) | 4 hours maximum | Disaster recovery testing | M1 |
| NFR-017 | Recovery Point Objective (RPO) | 1 hour maximum data loss | Backup and recovery testing | M1 |
| NFR-018 | Mean Time To Recovery (MTTR) | < 2 hours | Incident response metrics | M2 |
| NFR-019 | Service degradation handling | Graceful degradation during outages | Chaos engineering testing | M2 |
| NFR-020 | Data backup frequency | Real-time replication + daily backups | Backup monitoring | M1 |

## 5.5 Security Requirements
*Refer to **Section 7: Security & Compliance Requirements**  for detailed security controls.*

## 5.6 Usability Requirements

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-029 | User satisfaction score | 4.5+ stars on app stores | User feedback and ratings | M2 |
| NFR-030 | Task completion rate | 95% for critical user journeys | User experience analytics | M1 |
| NFR-031 | Learning curve | New users complete first task < 5 minutes | User onboarding metrics | M2 |
| NFR-032 | Error recovery | Clear error messages with action guidance | Error tracking and analysis | M1 |
| NFR-033 | Accessibility compliance | WCAG 2.1 AA compliance | Accessibility testing tools | M2 |
| NFR-034 | Multi-language support | Bengali and English localization | Localization testing | M1 |

## 5.7 Maintainability & Operability

| NFR ID | Requirement | Target | Measurement | Priority |
|--------|-------------|--------|-------------|----------|
| NFR-035 | Code coverage | 80% unit test coverage | Automated testing reports | M2 |
| NFR-036 | Deployment frequency | Daily deployments capability | CI/CD pipeline metrics | M2 |
| NFR-037 | Mean Time To Deploy | < 30 minutes for hotfixes | Deployment automation metrics | M2 |
| NFR-038 | Monitoring coverage | 100% critical path monitoring | Observability platform | M1 |
| NFR-039 | Log aggregation | Centralized logging for all services | Logging platform metrics | M1 |
| NFR-040 | Documentation currency | API documentation auto-generated | Documentation automation | M2 |

## 5.8 Compliance Requirements
*Refer to **Section 7: Security & Compliance Requirements** for detailed IDRA and BFIU compliance frameworks.*

