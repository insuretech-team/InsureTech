# B2B Events Proto Definitions

This directory contains all event definitions for the B2B domain.

## File Structure

- `core.proto` - All B2B domain events (organisation, members, departments, employees, purchase orders)

## Event Categories

### Organisation Events
- `OrganisationCreatedEvent` - New organisation created by SuperAdmin
- `OrganisationUpdatedEvent` - Organisation details updated
- `OrganisationStatusChangedEvent` - Organisation status changed
- `OrganisationApprovedEvent` - Organisation approved by SuperAdmin
- `OrganisationSuspendedEvent` - Organisation suspended

### OrgMember Events
- `OrgMemberAddedEvent` - User added to organisation
- `OrgMemberRemovedEvent` - User removed from organisation
- `OrgMemberRoleChangedEvent` - Member role changed
- `B2BAdminAssignedEvent` - B2B admin assigned to organisation

### Department Events
- `DepartmentCreatedEvent` - New department created
- `DepartmentUpdatedEvent` - Department updated
- `DepartmentDeletedEvent` - Department deleted

### Employee Events
- `EmployeeCreatedEvent` - New employee created
- `EmployeeUpdatedEvent` - Employee updated
- `EmployeeDeletedEvent` - Employee deleted
- `EmployeeStatusChangedEvent` - Employee status changed

### Purchase Order Events
- `PurchaseOrderCreatedEvent` - New purchase order created
- `PurchaseOrderStatusChangedEvent` - Purchase order status changed
- `PurchaseOrderApprovedEvent` - Purchase order approved
- `PurchaseOrderRejectedEvent` - Purchase order rejected

## Kafka Topics

All events are published to Kafka with the following topic naming convention:

- `b2b.organisation.created`
- `b2b.organisation.updated`
- `b2b.organisation.status_changed`
- `b2b.organisation.approved`
- `b2b.organisation.suspended`
- `b2b.org_member.added`
- `b2b.org_member.removed`
- `b2b.org_member.role_changed`
- `b2b.admin.assigned`
- `b2b.department.created`
- `b2b.department.updated`
- `b2b.department.deleted`
- `b2b.employee.created`
- `b2b.employee.updated`
- `b2b.employee.deleted`
- `b2b.employee.status_changed`
- `b2b.purchase_order.created`
- `b2b.purchase_order.status_changed`
- `b2b.purchase_order.approved`
- `b2b.purchase_order.rejected`

## Event Sourcing

All events are:
- Immutable
- Append-only
- Timestamped
- Include event_id for idempotency
- Include actor information (created_by, updated_by, etc.)

## Usage

Events are published by the B2B service and consumed by:
- AuthZ service (for permission management)
- Notification service (for alerts)
- Analytics service (for reporting)
- Audit service (for compliance)
