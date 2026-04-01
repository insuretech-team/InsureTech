-- =====================================================
-- Fix: widen VARCHAR columns in workflow_schema to fit proto enum name strings
--
-- Proto enum names are stored as their full string representation, e.g.:
--   WORKFLOW_STATUS_ACTIVE (21 chars)   → was VARCHAR(20) → needs VARCHAR(50)
--   WORKFLOW_TYPE_APPROVAL (22 chars)   → was VARCHAR(20) → needs VARCHAR(50)
--   INSTANCE_STATUS_IN_PROGRESS (27 chars) → needs VARCHAR(50)
--   WORKFLOW_TASK_TYPE_APPROVAL (27 chars) → needs VARCHAR(50)
--   WORKFLOW_TASK_STATUS_PENDING (28 chars) → needs VARCHAR(50)
--
-- All enum columns are widened to VARCHAR(50) to accommodate current and future
-- proto enum value name strings.
-- =====================================================

BEGIN;

-- ── workflow_definitions ──────────────────────────────────────────────────────

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_definitions'
                 AND column_name = 'type') THEN
        ALTER TABLE workflow_schema.workflow_definitions
            ALTER COLUMN type TYPE VARCHAR(50);
    END IF;
END $$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_definitions'
                 AND column_name = 'status') THEN
        ALTER TABLE workflow_schema.workflow_definitions
            ALTER COLUMN status TYPE VARCHAR(50);
    END IF;
END $$;

-- Update CHECK constraints to match widened columns
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_definitions'
                 AND constraint_name = 'chk_workflow_definitions_type') THEN
        ALTER TABLE workflow_schema.workflow_definitions
            DROP CONSTRAINT chk_workflow_definitions_type;
    END IF;

    ALTER TABLE workflow_schema.workflow_definitions
        ADD CONSTRAINT chk_workflow_definitions_type
        CHECK (type IN (
            'WORKFLOW_TYPE_UNSPECIFIED',
            'WORKFLOW_TYPE_APPROVAL',
            'WORKFLOW_TYPE_REVIEW',
            'WORKFLOW_TYPE_ESCALATION',
            'WORKFLOW_TYPE_NOTIFICATION'
        ));
END $$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_definitions'
                 AND constraint_name = 'chk_workflow_definitions_status') THEN
        ALTER TABLE workflow_schema.workflow_definitions
            DROP CONSTRAINT chk_workflow_definitions_status;
    END IF;

    ALTER TABLE workflow_schema.workflow_definitions
        ADD CONSTRAINT chk_workflow_definitions_status
        CHECK (status IN (
            'WORKFLOW_STATUS_UNSPECIFIED',
            'WORKFLOW_STATUS_DRAFT',
            'WORKFLOW_STATUS_ACTIVE',
            'WORKFLOW_STATUS_INACTIVE',
            'WORKFLOW_STATUS_ARCHIVED'
        ));
END $$;

-- ── workflow_instances ────────────────────────────────────────────────────────

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_instances'
                 AND column_name = 'status') THEN
        ALTER TABLE workflow_schema.workflow_instances
            ALTER COLUMN status TYPE VARCHAR(50);
    END IF;
END $$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_instances'
                 AND constraint_name = 'chk_workflow_instances_status') THEN
        ALTER TABLE workflow_schema.workflow_instances
            DROP CONSTRAINT chk_workflow_instances_status;
    END IF;

    ALTER TABLE workflow_schema.workflow_instances
        ADD CONSTRAINT chk_workflow_instances_status
        CHECK (status IN (
            'INSTANCE_STATUS_UNSPECIFIED',
            'INSTANCE_STATUS_PENDING',
            'INSTANCE_STATUS_IN_PROGRESS',
            'INSTANCE_STATUS_COMPLETED',
            'INSTANCE_STATUS_FAILED',
            'INSTANCE_STATUS_CANCELLED'
        ));
END $$;

-- ── workflow_tasks ────────────────────────────────────────────────────────────

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_tasks'
                 AND column_name = 'type') THEN
        ALTER TABLE workflow_schema.workflow_tasks
            ALTER COLUMN type TYPE VARCHAR(50);
    END IF;
END $$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_tasks'
                 AND column_name = 'status') THEN
        ALTER TABLE workflow_schema.workflow_tasks
            ALTER COLUMN status TYPE VARCHAR(50);
    END IF;
END $$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_tasks'
                 AND constraint_name = 'chk_workflow_tasks_type') THEN
        ALTER TABLE workflow_schema.workflow_tasks
            DROP CONSTRAINT chk_workflow_tasks_type;
    END IF;

    ALTER TABLE workflow_schema.workflow_tasks
        ADD CONSTRAINT chk_workflow_tasks_type
        CHECK (type IN (
            'WORKFLOW_TASK_TYPE_UNSPECIFIED',
            'WORKFLOW_TASK_TYPE_APPROVAL',
            'WORKFLOW_TASK_TYPE_REVIEW',
            'WORKFLOW_TASK_TYPE_NOTIFICATION',
            'WORKFLOW_TASK_TYPE_ACTION'
        ));
END $$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.table_constraints
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_tasks'
                 AND constraint_name = 'chk_workflow_tasks_status') THEN
        ALTER TABLE workflow_schema.workflow_tasks
            DROP CONSTRAINT chk_workflow_tasks_status;
    END IF;

    ALTER TABLE workflow_schema.workflow_tasks
        ADD CONSTRAINT chk_workflow_tasks_status
        CHECK (status IN (
            'WORKFLOW_TASK_STATUS_UNSPECIFIED',
            'WORKFLOW_TASK_STATUS_PENDING',
            'WORKFLOW_TASK_STATUS_IN_PROGRESS',
            'WORKFLOW_TASK_STATUS_COMPLETED',
            'WORKFLOW_TASK_STATUS_SKIPPED'
        ));
END $$;

-- ── workflow_configs ──────────────────────────────────────────────────────────

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'workflow_schema'
                 AND table_name = 'workflow_configs'
                 AND column_name = 'config_type') THEN
        ALTER TABLE workflow_schema.workflow_configs
            ALTER COLUMN config_type TYPE VARCHAR(60);
    END IF;
END $$;

COMMIT;
