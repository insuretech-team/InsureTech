package enforcer

// pg_adapter.go — Custom Casbin persist.Adapter backed by authz_schema.casbin_rules.
//
// WHY: The official gorm-adapter/v3 NewAdapterByDBWithCustomTable does not support
// schema-qualified table names (e.g. "authz_schema.casbin_rules"). Passing such a
// string causes the adapter to create a table literally named "authz_schema.casbin_rules"
// in the public schema. This adapter uses raw SQL via database/sql so every query
// explicitly targets authz_schema.casbin_rules — matching the proto source of truth.
//
// Casbin persist.Adapter interface (v3):
//   LoadPolicy(model model.Model) error
//   SavePolicy(model model.Model) error
//   AddPolicy(sec string, ptype string, rule []string) error
//   AddPolicies(sec string, ptype string, rules [][]string) error
//   RemovePolicy(sec string, ptype string, rule []string) error
//   RemovePolicies(sec string, ptype string, rules [][]string) error
//   RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"gorm.io/gorm"
)

// loadPolicyLine converts a ptype + rule slice into the CSV format that
// casbin/v3 persist.LoadPolicyLine expects: "p, sub, dom, obj, act, eft"
func loadPolicyLine(ptype string, rule []string, m model.Model) error {
	line := ptype + ", " + strings.Join(rule, ", ")
	return persist.LoadPolicyLine(line, m)
}

const casbinTable = "authz_schema.casbin_rules"

// pgAdapter implements casbin/v3 persist.BatchAdapter.
type pgAdapter struct {
	db *sql.DB
}

// newPGAdapter creates a new pgAdapter from a *gorm.DB.
// It extracts the underlying *sql.DB to use raw queries only.
func newPGAdapter(gdb *gorm.DB) (*pgAdapter, error) {
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("pg_adapter: get sql.DB: %w", err)
	}
	return &pgAdapter{db: sqlDB}, nil
}

// casbinRow mirrors a row in authz_schema.casbin_rules.
type casbinRow struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func rowToRule(r casbinRow) []string {
	rule := []string{r.V0}
	for _, v := range []string{r.V1, r.V2, r.V3, r.V4, r.V5} {
		if v == "" {
			break
		}
		rule = append(rule, v)
	}
	return rule
}

// ── persist.Adapter interface ─────────────────────────────────────────────────

// LoadPolicy loads all rules from authz_schema.casbin_rules into the Casbin model.
func (a *pgAdapter) LoadPolicy(m model.Model) error {
	rows, err := a.db.Query(
		`SELECT ptype, COALESCE(v0,''), COALESCE(v1,''), COALESCE(v2,''), COALESCE(v3,''), COALESCE(v4,''), COALESCE(v5,'') FROM ` + casbinTable + ` ORDER BY id`,
	)
	if err != nil {
		return fmt.Errorf("pg_adapter LoadPolicy: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r casbinRow
		if err := rows.Scan(&r.PType, &r.V0, &r.V1, &r.V2, &r.V3, &r.V4, &r.V5); err != nil {
			return fmt.Errorf("pg_adapter LoadPolicy scan: %w", err)
		}
		rule := rowToRule(r)
		if err := loadPolicyLine(r.PType, rule, m); err != nil {
			return fmt.Errorf("pg_adapter LoadPolicy line: %w", err)
		}
	}
	return rows.Err()
}

// SavePolicy replaces all rules with the current Casbin model state.
func (a *pgAdapter) SavePolicy(m model.Model) error {
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("pg_adapter SavePolicy begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM ` + casbinTable); err != nil {
		return fmt.Errorf("pg_adapter SavePolicy delete: %w", err)
	}

	var rows []casbinRow
	for ptype, assertions := range m["p"] {
		for _, rule := range assertions.Policy {
			rows = append(rows, ruleToRow(ptype, rule))
		}
	}
	for ptype, assertions := range m["g"] {
		for _, rule := range assertions.Policy {
			rows = append(rows, ruleToRow(ptype, rule))
		}
	}

	for _, r := range rows {
		if _, err := tx.Exec(
			`INSERT INTO `+casbinTable+` (ptype,v0,v1,v2,v3,v4,v5)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING`,
			r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5,
		); err != nil {
			return fmt.Errorf("pg_adapter SavePolicy insert: %w", err)
		}
	}
	return tx.Commit()
}

// AddPolicy inserts a single p or g rule.
func (a *pgAdapter) AddPolicy(sec, ptype string, rule []string) error {
	r := ruleToRow(ptype, rule)
	_, err := a.db.Exec(
		`INSERT INTO `+casbinTable+` (ptype,v0,v1,v2,v3,v4,v5)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING`,
		r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5,
	)
	if err != nil {
		return fmt.Errorf("pg_adapter AddPolicy: %w", err)
	}
	return nil
}

// AddPolicies inserts multiple rules in a single transaction.
func (a *pgAdapter) AddPolicies(sec, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("pg_adapter AddPolicies begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(
		`INSERT INTO ` + casbinTable + ` (ptype,v0,v1,v2,v3,v4,v5)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)
		 ON CONFLICT ON CONSTRAINT uq_casbin_rules_tuple DO NOTHING`,
	)
	if err != nil {
		return fmt.Errorf("pg_adapter AddPolicies prepare: %w", err)
	}
	defer stmt.Close()

	for _, rule := range rules {
		r := ruleToRow(ptype, rule)
		if _, err := stmt.Exec(r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5); err != nil {
			return fmt.Errorf("pg_adapter AddPolicies exec: %w", err)
		}
	}
	return tx.Commit()
}

// RemovePolicy deletes a single rule matching all provided fields.
func (a *pgAdapter) RemovePolicy(sec, ptype string, rule []string) error {
	return a.RemovePolicies(sec, ptype, [][]string{rule})
}

// RemovePolicies deletes multiple rules in a single transaction.
func (a *pgAdapter) RemovePolicies(sec, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("pg_adapter RemovePolicies begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, rule := range rules {
		r := ruleToRow(ptype, rule)
		if _, err := tx.Exec(
			`DELETE FROM `+casbinTable+`
			 WHERE ptype=$1 AND v0=$2 AND v1=$3 AND v2=$4 AND v3=$5 AND v4=$6 AND v5=$7`,
			r.PType, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5,
		); err != nil {
			return fmt.Errorf("pg_adapter RemovePolicies: %w", err)
		}
	}
	return tx.Commit()
}

// RemoveFilteredPolicy deletes rules matching ptype and field values starting at fieldIndex.
func (a *pgAdapter) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	cols := []string{"v0", "v1", "v2", "v3", "v4", "v5"}
	var conditions []string
	var args []interface{}
	args = append(args, ptype)
	conditions = append(conditions, "ptype = $1")

	for i, v := range fieldValues {
		if v == "" {
			continue
		}
		col := cols[fieldIndex+i]
		args = append(args, v)
		conditions = append(conditions, fmt.Sprintf("%s = $%d", col, len(args)))
	}

	query := `DELETE FROM ` + casbinTable + ` WHERE ` + strings.Join(conditions, " AND ")
	if _, err := a.db.Exec(query, args...); err != nil {
		return fmt.Errorf("pg_adapter RemoveFilteredPolicy: %w", err)
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func ruleToRow(ptype string, rule []string) casbinRow {
	r := casbinRow{PType: ptype}
	vals := []*string{&r.V0, &r.V1, &r.V2, &r.V3, &r.V4, &r.V5}
	for i, v := range rule {
		if i >= len(vals) {
			break
		}
		*vals[i] = v
	}
	return r
}

// Compile-time interface checks.
var _ persist.Adapter = (*pgAdapter)(nil)
var _ persist.BatchAdapter = (*pgAdapter)(nil)

// Ensure errors package is used (for future error wrapping in this file).
var _ = errors.New
