# Improved Moodle Automated Course Backup

Drop-in replacement for Moodle's ./admin/cli/automated_backups.php or Moodle's regular cron running automated backups.

Replaces the regular automated backups task with more efficient, concurrent, controlled version for large sites.

It operates by the definitions from Moodle's settings for automated backups from the database.
