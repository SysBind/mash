# Improved Moodle Automated Course Backup

Drop-in replacement for Moodle's ./admin/cli/automated_backups.php or Moodle's regular cron running automated backups.

Replaces the regular automated backups task with more efficient, concurrent, controlled version for large sites.

It operates by the definitions from Moodle's settings for automated backups from the database.


## moodle go module

This will serve as a POC for the moodle/ go package

### Unit Testing

Copy moodle/testdata/config-dist.php to moodle/testdata/config.php,
Adjust Parameters, Run clean installation of moodle (>=3.8),
This config.php need to contain $CFG->dirroot pointing to moodle's code
Than ```go test ./...```
