# config-fs
Config fs is a small application for holding mongodb collection data in file system

**Idea**

Config FS is an cli application to store configuration in file system and to sync it with mongodb

It is useful to store configurations in git repo and sync it to mongodb collections

folder <=> collection

file <=> record inside collection

read folder and store data in database 

config-fs read --database=ug --collection "websites"

read collection website from mongodb database named ug and store data in fs

config-fs read -d ug -c "website" --connection="mongodb://10.0.1.77:27017" /Users/taleh/Projects/config-fs/test

read fs and store records in database ug

config-fs read -d ug -c "website" --connection="mongodb://10.0.1.77:27017" /Users/taleh/Projects/config-fs/test

