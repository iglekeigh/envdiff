// Package audit provides a lightweight audit log for tracking changes made
// to environment variables during diff, reconcile, merge, and redact operations.
//
// Each operation that modifies or inspects sensitive keys can record an Entry
// describing what changed, when, and from which source file or label.
//
// Usage:
//
//	log := audit.New("production")
//	log.Record("DB_PASSWORD", audit.ActionRedacted, "", "")
//	for _, e := range log.Entries() {
//		fmt.Println(e)
//	}
package audit
