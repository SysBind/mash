## Mash - Alternative Moodle Shell Written in GO

## Improved Moodle Automated Course Backup

Copy moodle/testdata/automated_backup_single.php under Moodle's admin/cli.
(helper scripts will be embedded in the Mash binary in the future)

From Moodle's root:
```mash autobackup```

Drop-in replacement for Moodle's ./admin/cli/automated_backups.php or Moodle's regular cron running automated backups.

Replaces the regular automated backups task with more efficient, concurrent, controlled version for large sites.

It operates by the definitions from Moodle's settings for automated backups from the database.


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
