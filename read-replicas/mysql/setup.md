docker exec -it mysql-primary mysql -uroot -p
Enter password:
Welcome to the MySQL monitor. Commands end with ; or \g.
Your MySQL connection id is 9
Server version: 8.0.44 MySQL Community Server - GPL

Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> CREATE USER 'replication_user'@'%' IDENTIFIED WITH mysql_native_password BY 'replication_password';
Query OK, 0 rows affected (0.01 sec)

mysql> GRANT REPLICATION SLAVE ON _._ TO 'replication_user'@'%';
Query OK, 0 rows affected (0.00 sec)

mysql> FLUSH PRIVILEGES;
Query OK, 0 rows affected (0.01 sec)

mysql> SHOW MASTER STATUS;
+------------------+----------+--------------+------------------+-------------------+
| File | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000003 | 851 | | | |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

docker exec -it mysql-replica mysql -uroot -preplica_password
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor. Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.0.44 MySQL Community Server - GPL

Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> CHANGE REPLICATION SOURCE TO
-> SOURCE_HOST='mysql-primary',
-> SOURCE_USER='replication_user',
-> SOURCE_PASSWORD='replication_password',
-> SOURCE_LOG_FILE='mysql-bin.000003',
-> SOURCE_LOG_POS=851;

Query OK, 0 rows affected, 2 warnings (0.48 sec)

mysql>
mysql> START REPLICA;
Query OK, 0 rows affected (0.03 sec)

mysql> SHOW REPLICA STATUS\G \***\*\*\*\*\*\*\***\*\*\*\***\*\*\*\*\*\*\*** 1. row \***\*\*\*\*\*\*\***\*\*\*\***\*\*\*\*\*\*\***
Replica_IO_State: Checking source version
Source_Host: mysql-primary
Source_User: replication_user
Source_Port: 3306
Connect_Retry: 60
Source_Log_File: mysql-bin.000003
Read_Source_Log_Pos: 851
Relay_Log_File: relay-log-bin.000001
Relay_Log_Pos: 4
Relay_Source_Log_File: mysql-bin.000003
Replica_IO_Running: Yes
Replica_SQL_Running: Yes
Replicate_Do_DB:
Replicate_Ignore_DB:
Replicate_Do_Table:
Replicate_Ignore_Table:
Replicate_Wild_Do_Table:
Replicate_Wild_Ignore_Table:
Last_Errno: 0
Last_Error:
Skip_Counter: 0
Exec_Source_Log_Pos: 851
Relay_Log_Space: 157
Until_Condition: None
Until_Log_File:
Until_Log_Pos: 0
Source_SSL_Allowed: No
Source_SSL_CA_File:
Source_SSL_CA_Path:
Source_SSL_Cert:
Source_SSL_Cipher:
Source_SSL_Key:
Seconds_Behind_Source: 0
Source_SSL_Verify_Server_Cert: No
Last_IO_Errno: 0
Last_IO_Error:
Last_SQL_Errno: 0
Last_SQL_Error:
Replicate_Ignore_Server_Ids:
Source_Server_Id: 0
Source_UUID:
Source_Info_File: mysql.slave_master_info
SQL_Delay: 0
SQL_Remaining_Delay: NULL
Replica_SQL_Running_State: Replica has read all relay log; waiting for more updates
Source_Retry_Count: 86400
Source_Bind:
Last_IO_Error_Timestamp:
Last_SQL_Error_Timestamp:
Source_SSL_Crl:
Source_SSL_Crlpath:
Retrieved_Gtid_Set:
Executed_Gtid_Set:
Auto_Position: 0
Replicate_Rewrite_DB:
Channel_Name:
Source_TLS_Version:
Source_public_key_path:
Get_Source_public_key: 0
Network_Namespace:
1 row in set (0.00 sec)

https://claude.ai/share/8701b10c-94b5-4797-a068-a04ba845c0bc
https://claude.ai/share/8701b10c-94b5-4797-a068-a04ba845c0bc
