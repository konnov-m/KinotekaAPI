FROM postgres:15

COPY migrate/create_db.sql migrate/insert.sql migrate/migrate.sh /migrate/
RUN chmod +x migrate/migrate.sh

ENTRYPOINT ["/migrate/migrate.sh"]