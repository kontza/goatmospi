#
# Run with 'changeLogSync' on an existing table.
# Run with 'migrate' when the have been changes.
#
liquibase \
    --classpath=./postgresql-42.1.1.jar \
    --driver=org.postgresql.Driver \
    --changeLogFile=goatmospi-changelog.xml \
    --url="jdbc:postgresql://localhost/atmospi" \
    --username=atmospi \
    --password=atmospi \
    $argv[1]