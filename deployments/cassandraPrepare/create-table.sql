CREATE keyspace finch
WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};
CREATE TABLE finch.urls (
        id text primary key ,
        url text
);
