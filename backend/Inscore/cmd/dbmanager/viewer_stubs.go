package main

import (
	"fmt"

	"github.com/newage-saint/insuretech/backend/inscore/db/ops"
)

// RunInteractiveTableViewer launches the interactive table viewer TUI (stub)
func RunInteractiveTableViewer() error {
	fmt.Println("Interactive table viewer not yet implemented")
	return nil
}

// RunSyncReportViewer launches the sync report viewer TUI (stub)
func RunSyncReportViewer(statuses []ops.TableStatus) error {
	fmt.Println("Sync report viewer not yet implemented")
	fmt.Printf("Would display %d tables\n", len(statuses))
	return nil
}

