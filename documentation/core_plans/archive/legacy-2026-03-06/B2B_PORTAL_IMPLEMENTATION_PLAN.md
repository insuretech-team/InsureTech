# B2B Portal Implementation Plan
## Comprehensive List of Unimplemented Features & Integration Roadmap

**Generated:** March 2026  
**Status:** All services running, core features implemented  
**Next Phase:** Advanced features, backend integrations, policy/document management

---

## ✅ **IMPLEMENTED FEATURES** (Working & Tested)

| Feature | Status | Backend Integration | Notes |
|---------|--------|-------------------|-------|
| **Dashboard** | ✅ Complete | B2B + stats/activity APIs | Live KPI stats, recent activity feed |
| **Organizations** | ✅ Complete | B2B service | CRUD, approve, member management |
| **Departments** | ✅ Complete | B2B service | Org-scoped CRUD, Super Admin vs B2B Admin |
| **Employees** | ✅ Complete | B2B service | Full record edit, department filtering |
| **Purchase Orders** | ✅ Complete | B2B service | Create, view, delete, catalog integration |
| **Authentication** | ✅ Complete | AuthN service (:50060) | JWT, mobile normalization, role resolution |
| **Authorization** | ✅ Complete | AuthZ service (:50070) | Casbin PERM model, org-scoped permissions |
| **Team Management** | ✅ Complete | B2B service | HR/Viewer role management for B2B Admin |

---

## 🔴 **UNIMPLEMENTED FEATURES** (Mock Data Only)

### **Tier 1: Critical Business Features** (Start with these)

| Feature | Current State | Required Integration | Priority |
|---------|---------------|----------------------|----------|
| **Insurance Plans** | 📊 Mock data UI | **Insurance Service** (:50XX) | 🔥 **HIGH** |
| **Policies** | 📊 Mock data UI | **Insurance Service** + **DocGen** | 🔥 **HIGH** |
| **Billing & Invoices** | 📊 Mock data UI | **Payment Service** + **Storage** | 🔥 **HIGH** |
| **Claims** | 📊 Mock data UI | **Insurance Service** + **DocGen** + **Storage** | 🔥 **HIGH** |
| **Payments** | 📊 Mock data UI | **Payment Service** | 🔥 **HIGH** |
| **Settings → Organization Form** | 📊 Mock data UI | B2B service (wiring needed) | 🟡 **MEDIUM** |
| **Settings → Workflow Rules** | 📊 Mock data UI | New workflow API | 🟡 **MEDIUM** |
| **Settings → Notifications** | 📊 Mock data UI | **Notification Service** | 🟡 **MEDIUM** |

### **Tier 2: Advanced Features** (After Tier 1)

| Feature | Current State | Required Integration | Priority |
|---------|---------------|----------------------|----------|
| **Document Generation** | ❌ Not implemented | **DocGen Service** (:50XX) | 🟡 **MEDIUM** |
| **File Upload/Storage** | ❌ Not implemented | **Storage Service** | 🟡 **MEDIUM** |
| **Audit Trails** | ❌ Not implemented | **Audit Service** | 🟢 **LOW** |
| **Analytics Dashboard** | ❌ Not implemented | **Analytics Service** | 🟢 **LOW** |
| **KYC Verification** | ❌ Not implemented | **KYC Service** | 🟢 **LOW** |
| **Fraud Detection** | ❌ Not implemented | **Fraud Service** | 🟢 **LOW** |

### **Tier 3: Missing Pages** (No routes exist)

| Page | Route | Current State |
|------|-------|---------------|
| **Quotations** | `/quotations` | ❌ Route missing, mockups exist in `b2b-dashboard-client.ts` |
| **Reports** | `/reports` | ❌ Not mentioned in navigation |
| **Profile** | `/profile` | ❌ Route missing |
| **Help/Support** | `/support` | ❌ Route missing |

---

## 🏗️ **BACKEND SERVICES STATUS**

| Service | Port | Status | Notes |
|---------|------|---------|-------|
| **B2B** | :50112 | ✅ Running | Core org/employee/PO management |
| **AuthN** | :50060 | ✅ Running | User authentication |
| **AuthZ** | :50070 | ✅ Running | Casbin permissions |
| **Gateway** | :8080 | ✅ Running | REST API gateway |
| **Kafka** | :9092 | ✅ Running | Event streaming |
| **DocGen** | :50XX | 🟡 Available | PDF/document generation |
| **Insurance** | :50XX | 🟡 Available | Policy/claims management |
| **Media** | :50XX | 🟡 Available | File upload/storage |
| **Fraud** | :50XX | 🟡 Available | Fraud detection |
| **Payment** | ❌ | ❌ Missing | **NEED TO IMPLEMENT** |
| **Storage** | ❌ | ❌ Missing | **NEED TO IMPLEMENT** |
| **Notification** | ❌ | ❌ Missing | **NEED TO IMPLEMENT** |
| **Analytics** | ❌ | ❌ Missing | **NEED TO IMPLEMENT** |
| **Audit** | ❌ | ❌ Missing | **NEED TO IMPLEMENT** |
| **KYC** | ❌ | ❌ Missing | **NEED TO IMPLEMENT** |

---

## 📋 **IMPLEMENTATION PLAN**

### **Phase 1: Core Business Features** (4-6 weeks)

#### **Week 1-2: Insurance Management**
1. **Wire Insurance Plans page** (`/insurance-plans`)
   - Create `/api/insurance-plans` BFF routes
   - Integrate with **Insurance Service**
   - Replace mock data with live plan catalog
   - Add plan selection/enrollment flows

2. **Wire Policies page** (`/policies`)
   - Create `/api/policies` BFF routes
   - Integrate with **Insurance Service**
   - Policy CRUD operations
   - Policy document generation via **DocGen Service**

#### **Week 3-4: Payment & Billing**
3. **Implement Payment Service**
   - Create new microservice: `backend/inscore/microservices/payment/`
   - gRPC interface: payment methods, transactions, billing
   - Database schema: payments, invoices, billing cycles

4. **Wire Payments & Billing pages**
   - Create `/api/payments`, `/api/billing-invoices` BFF routes
   - Replace mock data with live payment/invoice data
   - Payment processing workflows
   - Invoice generation and management

#### **Week 5-6: Claims Management**
5. **Wire Claims page** (`/claims`)
   - Create `/api/claims` BFF routes
   - Integrate with **Insurance Service**
   - Claims submission, tracking, approval workflows
   - Document upload for claims evidence

### **Phase 2: Storage & Document Management** (3-4 weeks)

#### **Week 7-8: Storage Infrastructure**
6. **Implement Storage Service**
   - Create new microservice: `backend/inscore/microservices/storage/`
   - File upload/download APIs
   - S3/MinIO integration for document storage
   - File metadata and versioning

7. **Integrate DocGen Service**
   - Policy documents, certificates, invoices
   - PDF generation for claims, reports
   - Template management

#### **Week 9-10: Document Workflows**
8. **Add file upload components**
   - Document upload forms for claims, employee records
   - File preview and download functionality
   - Document approval workflows

### **Phase 3: Settings & Workflow** (2-3 weeks)

#### **Week 11-12: Settings Management**
9. **Wire Settings page tabs**
   - Organization Profile → update via B2B service
   - Workflow Rules → implement workflow API
   - Notification Preferences → integrate Notification Service

10. **Implement Notification Service**
    - Create new microservice: `backend/inscore/microservices/notification/`
    - Email/SMS notification templates
    - Event-driven notifications (Kafka integration)

### **Phase 4: Analytics & Audit** (3-4 weeks)

#### **Week 13-14: Analytics**
11. **Implement Analytics Service**
    - Create new microservice: `backend/inscore/microservices/analytics/`
    - Business intelligence dashboards
    - Report generation

12. **Enhanced Dashboard**
    - Advanced KPI charts
    - Trend analysis
    - Export functionality

#### **Week 15-16: Audit & Compliance**
13. **Implement Audit Service**
    - Create new microservice: `backend/inscore/microservices/audit/`
    - Comprehensive audit trails
    - Compliance reporting

14. **Additional Features**
    - KYC integration
    - Fraud detection alerts
    - Advanced user management

---

## 🔧 **TECHNICAL IMPLEMENTATION DETAILS**

### **New BFF API Routes Needed**

```typescript
// Insurance Management
GET/POST /api/insurance-plans
GET/PATCH/DELETE /api/insurance-plans/{id}
POST /api/insurance-plans/{id}/enroll

// Policy Management  
GET/POST /api/policies
GET/PATCH/DELETE /api/policies/{id}
POST /api/policies/{id}/renew
GET /api/policies/{id}/documents
POST /api/policies/{id}/generate-certificate

// Payment & Billing
GET/POST /api/payments
GET /api/payments/{id}
POST /api/payments/{id}/process
GET/POST /api/billing-invoices
GET /api/billing-invoices/{id}
POST /api/billing-invoices/{id}/pay

// Claims Management
GET/POST /api/claims
GET/PATCH/DELETE /api/claims/{id}
POST /api/claims/{id}/documents
POST /api/claims/{id}/approve
POST /api/claims/{id}/reject

// Document Management
POST /api/documents/upload
GET /api/documents/{id}
DELETE /api/documents/{id}
GET /api/documents/{id}/download

// Workflow & Settings
GET/POST /api/workflows
GET/PATCH /api/workflows/{id}
GET/PATCH /api/settings/organization
GET/PATCH /api/settings/notifications

// Analytics & Reports
GET /api/analytics/dashboard
GET /api/analytics/reports
POST /api/reports/generate
```

### **New Backend Services Architecture**

```
Payment Service (:50130)
├── PaymentController
├── InvoiceController  
├── BillingController
└── PaymentMethodController

Storage Service (:50140)
├── FileController
├── DocumentController
└── TemplateController

Notification Service (:50150)
├── NotificationController
├── TemplateController
└── PreferenceController

Analytics Service (:50160)
├── DashboardController
├── ReportController
└── MetricsController

Audit Service (:50170)
├── AuditLogController
├── ComplianceController
└── TrailController
```

### **Database Schema Extensions**

```sql
-- Payment Management
CREATE TABLE payment_schema.payment_methods (
    method_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    type VARCHAR(50), -- BANK_TRANSFER, CARD, BKASH, NAGAD
    details JSONB,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE payment_schema.invoices (
    invoice_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    amount DECIMAL(12,2),
    due_date DATE,
    status VARCHAR(20),
    items JSONB
);

-- Document Storage
CREATE TABLE storage_schema.documents (
    document_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    file_path VARCHAR(500),
    file_size BIGINT,
    mime_type VARCHAR(100),
    metadata JSONB
);

-- Workflow Management
CREATE TABLE workflow_schema.workflows (
    workflow_id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    name VARCHAR(200),
    steps JSONB,
    is_active BOOLEAN
);
```

---

## 🎯 **SUCCESS CRITERIA**

### **Phase 1 Completion:**
- [ ] Insurance plans browsing and enrollment
- [ ] Policy lifecycle management  
- [ ] Payment processing workflows
- [ ] Claims submission and tracking
- [ ] All mock data replaced with live APIs

### **Phase 2 Completion:**
- [ ] Document upload/download functionality
- [ ] PDF generation for policies and claims
- [ ] File storage and management

### **Phase 3 Completion:**
- [ ] Complete settings management
- [ ] Workflow automation
- [ ] Notification system

### **Phase 4 Completion:**
- [ ] Business intelligence dashboards
- [ ] Comprehensive audit trails
- [ ] Compliance reporting

### **Overall Success Metrics:**
- [ ] Zero mock data remaining
- [ ] All navigation items functional
- [ ] Complete insurance business workflow
- [ ] Document management integration
- [ ] Advanced analytics and reporting
- [ ] Full audit compliance

---

## 🚀 **NEXT IMMEDIATE ACTIONS**

1. **Prioritize Tier 1 features** — start with Insurance Plans integration
2. **Set up development environment** for new microservices
3. **Design database schemas** for Payment and Storage services
4. **Create API specifications** for new BFF routes
5. **Plan integration testing** strategy across all services

---

*This plan provides a clear roadmap to transform the B2B portal from a functional MVP with some mock data into a complete enterprise-grade insurance management platform.*