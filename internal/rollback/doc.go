// Package rollback provides functionality to restore .env files
// from previously created backups. It lists available backups for
// a given env file and copies the selected backup back to the
// original path, enabling quick recovery from unintended changes.
package rollback
