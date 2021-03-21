## Mash - Alternative Moodle Shell Written in GO

## Improved Moodle Automated Course Backup

```mash autobackup```

Drop-in replacement for Moodle's ./admin/cli/automated_backups.php or Moodle's regular cron running automated backups.

Replaces the regular automated backups task with more efficient, concurrent, controlled version for large sites.

It operates by the definitions from Moodle's settings for automated backups from the database.

_Respected Settings_
- backup_auto_active
- backup_auto_max_kept
- backup_auto_storage
- backup_auto_destination
- backup_auto_skip_modif_prev
In addition, if [https://docs.moodle.org/310/en/Context_freezing](contextlocking) is true,
Frozen courses and courses under frozen categories will be skipped.


## moodle go module

This will serve as a POC for the moodle/ go package

### Unit Testing

Copy moodle/testdata/config-dist.php to moodle source direcotry,
Adjust Parameters, Run clean installation of moodle (>=3.8).

Than symlink to moodle/testdata/config.php,
This config.php need to contain $CFG->dirroot pointing to moodle's code

Than ```go test ./...```

Some tests will need to create multiple backups,
So it will sleep 1 minute between backups
(Backups are timestamped to minute, so if we dont sleep it will override each-other)
