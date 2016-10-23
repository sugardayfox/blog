Id: 1300
Title: Basics of mysql
Tags: sql,mysql
Date: 2006-01-14T16:00:00-08:00
Format: Markdown
Deleted: yes
--------------
**Start the console**: `mysql`

**List all databases**: `show databases;`

**Switch to a database**: `use $dbname;`

**List tables in a current database**: `show tables;`

**List columns in a given table**: `describe $tablename;`

**Setting permissions to a database**: `GRANT SELECT ON $database.* TO '$user_name'@'%' IDENTIFIED BY '$password';`

**Deleting rows**: `delete from $db.$table where $condition;`

**Selects**: `select count(*) from $db.$table;`

**Altering a table**:
```sql
ALTER TABLE get_cookie_log ADD COLUMN log_id INT(10) NOT NULL auto_increment, ADD PRIMARY KEY(log_id);`
ALTER TABLE get_cookie_log ADD COLUMN log_id INT(10) NOT NULL auto_increment, ADD INDEX(log_id);`
ALTER TABLE verify_reg_code_log ADD COLUMN log_id INT(10) NOT NULL auto_increment, ADD INDEX(log_id);`
```

**Backup database**: `mysqldump $database [[-all] | [$table]]s >$file-name.sql`

**Import database**:
```
mysql -e "DROP DATABASE $db; CREATE DATABASE $db;"
mysql $databasename <$file-name.sql
```
