FROM postgres:latest
COPY init.sql /docker-entrypoint-initdb.d/

# Make port available to the world outside this container
EXPOSE 5432