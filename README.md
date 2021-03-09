# Improved Moodle Automated Course Backup

Drop-in replacement for Moodle's ./admin/cli/automated_backups.php or Moodle's regular cron running automated backups.

Replaces the regular automated backups task with more efficient, concurrent, controlled version for large sites.

It operates by the definitions from Moodle's settings for automated backups from the database.


## moodle go module

This will serve as a POC for the moodle/ go package

### Unit Testing

Replace the moodle/testdata/config.php with config.php with clean installed moodle,
This config.php need to contain $CFG->dirroot pointing to moodle's code
(However, don't commit this file back..)
