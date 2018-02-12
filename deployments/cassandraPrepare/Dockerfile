FROM cassandra:3.11.1

ADD ./create-table.sql /tmp/

CMD [ "cqlsh","-f","/tmp/create-table.sql","cassandra"]