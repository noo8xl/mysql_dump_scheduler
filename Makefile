
# Create a database dump to backup
.PHONY: run-dump
run-dump:
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	mkdir -p ${DB_DUMP_DIR}; \
	mysqldump -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) > ${DB_DUMP_DIR}/$(DB_NAME).sql; \
	echo "Database dump created at ${DB_DUMP_DIR}/$(DB_NAME)_$${timestamp}.sql"	