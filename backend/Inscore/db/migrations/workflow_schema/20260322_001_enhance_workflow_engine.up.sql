-- =====================================================
-- Enhancement: workflow_schema — workflow engine tables
-- Phase 2 migration (proto-first system)
-- 
-- Tables are auto-created by Phase 1 proto-driven engine from:
--   proto/insuretech/workflow/entity/v1/workflow_definition.proto  (migration_order: 59)
--   proto/insuretech/workflow/entity/v1/workflow_instance.proto    (migration_order: 60)
--   proto/insuretech/workflow/entity/v1/workflow_task.proto        (migration_order: 61)
--   proto/insuretech/workflow/entity/v1/workflow_config.proto      (migration_order: 96)
--
-- This file adds only what the proto annotation system cannot express:
--   1. Composite indexes (multi-column — proto index only supports single column)
--   2. Partial indexes (WHERE clause — proto has no partial index support)
--   3. CHECK constraints on free-text columns (decision APPROVED/REJECTED/RETURNED)
--   4. updated_at auto-trigger (proto audit_fields adds the column but not the trigger)
--   5. Rich COMMENT ON COLUMN annotations for documentation
-- =====================================================

BEGIN;

-- ── 1. COMPOSITE & PARTIAL INDEXES ───────────────────────────────────────────
-- These cannot be expressed in proto column annotations

-- workflow_definitions: fast lookup of active definitions by entity type
-- Used by: WorkflowTriggerConsumer.handleMessage() and ListDefinitions()
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_definitions'
    ) THEN
        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_definitions_entity_status
            ON workflow_schema.workflow_definitions (entity_type, status)
            WHERE status = 'WORKFLOW_STATUS_ACTIVE'
        $idx$;

        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_definitions_name_trgm
            ON workflow_schema.workflow_definitions USING gin (name gin_trgm_ops)
        $idx$;
    END IF;
END $$;

-- workflow_instances: composite for "find all active instances for entity"
-- Used by: ListInstancesByEntity() — the workflow history query
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_instances'
    ) THEN
        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_instances_entity_status
            ON workflow_schema.workflow_instances (entity_type, entity_id, status)
            WHERE status IN ('INSTANCE_STATUS_PENDING', 'INSTANCE_STATUS_IN_PROGRESS')
        $idx$;

        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_instances_started_at
            ON workflow_schema.workflow_instances (started_at DESC)
        $idx$;

        -- Note: correlation_id is not in the current proto definition
        -- Add here if/when the proto is extended with that field
    END IF;
END $$;

-- workflow_tasks: the critical "my task inbox" composite index
-- Used by: ListTasksByAssignee() — the most frequent portal query
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_tasks'
    ) THEN
        -- Primary inbox index: assigned_to + pending status + due_date ordering
        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_tasks_inbox
            ON workflow_schema.workflow_tasks (assigned_to, status, due_date ASC NULLS LAST)
            WHERE assigned_to IS NOT NULL
              AND status IN ('WORKFLOW_TASK_STATUS_PENDING', 'WORKFLOW_TASK_STATUS_IN_PROGRESS')
        $idx$;

        -- Overdue tasks index: find tasks past due_date
        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_tasks_overdue
            ON workflow_schema.workflow_tasks (due_date ASC)
            WHERE due_date IS NOT NULL
              AND status IN ('WORKFLOW_TASK_STATUS_PENDING', 'WORKFLOW_TASK_STATUS_IN_PROGRESS')
        $idx$;
    END IF;
END $$;

-- workflow_configs: B2B per-business enabled configs
DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.tables
        WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_configs'
    ) THEN
        EXECUTE $idx$
            CREATE INDEX IF NOT EXISTS idx_wf_configs_business_enabled
            ON workflow_schema.workflow_configs (business_id, config_type)
            WHERE is_enabled = TRUE
        $idx$;
    END IF;
END $$;

-- ── 2. CHECK CONSTRAINTS ──────────────────────────────────────────────────────
-- Proto stores enum name strings in VARCHAR — add DB-level validation
-- for the decision column which is a free-text field (not a proto enum)

DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'workflow_schema'
          AND table_name = 'workflow_tasks'
          AND column_name = 'decision'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_schema = 'workflow_schema'
          AND table_name = 'workflow_tasks'
          AND constraint_name = 'chk_wf_tasks_decision'
    ) THEN
        ALTER TABLE workflow_schema.workflow_tasks
            ADD CONSTRAINT chk_wf_tasks_decision
            CHECK (decision IS NULL OR decision IN ('APPROVED', 'REJECTED', 'RETURNED'));
    END IF;
END $$;

-- ── 3. AUTO-UPDATE TRIGGERS ───────────────────────────────────────────────────
-- Proto audit_fields: true adds created_at + updated_at columns
-- but does NOT create the trigger — we add it here

DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'workflow_schema'
          AND table_name = 'workflow_definitions'
          AND column_name = 'updated_at'
    ) THEN
        EXECUTE $fn$
            CREATE OR REPLACE FUNCTION workflow_schema.trg_workflow_definitions_updated_at()
            RETURNS TRIGGER LANGUAGE plpgsql AS $body$
            BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
            $body$
        $fn$;

        DROP TRIGGER IF EXISTS trg_workflow_definitions_update
            ON workflow_schema.workflow_definitions;

        EXECUTE $tr$
            CREATE TRIGGER trg_workflow_definitions_update
            BEFORE UPDATE ON workflow_schema.workflow_definitions
            FOR EACH ROW EXECUTE FUNCTION workflow_schema.trg_workflow_definitions_updated_at()
        $tr$;
    END IF;
END $$;

DO $$ BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'workflow_schema'
          AND table_name = 'workflow_configs'
          AND column_name = 'updated_at'
    ) THEN
        EXECUTE $fn$
            CREATE OR REPLACE FUNCTION workflow_schema.trg_workflow_configs_updated_at()
            RETURNS TRIGGER LANGUAGE plpgsql AS $body$
            BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
            $body$
        $fn$;

        DROP TRIGGER IF EXISTS trg_workflow_configs_update
            ON workflow_schema.workflow_configs;

        EXECUTE $tr$
            CREATE TRIGGER trg_workflow_configs_update
            BEFORE UPDATE ON workflow_schema.workflow_configs
            FOR EACH ROW EXECUTE FUNCTION workflow_schema.trg_workflow_configs_updated_at()
        $tr$;
    END IF;
END $$;

-- ── 4. RICH COLUMN COMMENTS ───────────────────────────────────────────────────
-- Document runtime semantics that proto comments cannot fully capture

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables
               WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_definitions') THEN

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_definitions.steps IS
            'JSON array of step templates: [{name, type, assign_role, assign_to, due_hours, order}]. '
            'Parsed by Go workflow-engine at instance start to create WorkflowTask rows.' $c$;

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_definitions.conditions IS
            'JSON: {fail_fast_on_rejection, require_all_approvals, auto_approve_after_hours, escalate_to_role, custom_rules}. '
            'Evaluated by advanceInstance() to determine workflow completion outcome.' $c$;

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_definitions.entity_type IS
            'Domain entity this workflow applies to: CLAIM, ENDORSEMENT, REFUND, UNDERWRITING, POLICY, QUOTATION.' $c$;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables
               WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_instances') THEN

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_instances.entity_id IS
            'Polymorphic FK — references claim_id, endorsement_id, refund_id, policy_id etc. '
            'Resolve via entity_type to determine which schema/table to join.' $c$;

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_instances.context IS
            'Caller-supplied context JSON passed at StartWorkflow time. '
            'Used to carry claim amounts, customer details, or any domain metadata.' $c$;

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_instances.current_step IS
            'Name of the currently active step (from workflow_definitions.steps[].name). '
            'Updated by advanceInstance() as each task is completed.' $c$;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables
               WHERE table_schema = 'workflow_schema' AND table_name = 'workflow_tasks') THEN

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_tasks.decision IS
            'Set when task is completed: APPROVED, REJECTED, or RETURNED. '
            'Drives the advanceInstance() outcome logic in the Go workflow-engine.' $c$;

        EXECUTE $c$ COMMENT ON COLUMN workflow_schema.workflow_tasks.assigned_to IS
            'NULL until role is resolved to a specific user UUID. '
            'Role resolution happens at task creation time or via a role assignment service.' $c$;
    END IF;
END $$;

COMMIT;
