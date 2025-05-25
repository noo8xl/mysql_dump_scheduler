
# Create a database dump to backup
.PHONY: run-dump
run-dump:
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	mkdir -p ${DB_DUMP_DIR}; \
	mysqldump -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) > ${DB_DUMP_DIR}/$(DB_NAME).sql; \
	echo "Database dump created at ${DB_DUMP_DIR}/$(DB_NAME)_$${timestamp}.sql"	

gen-doc:
	@echo "Generating comprehensive documentation for mysql_dupm_scheduler..."
	@mkdir -p docs
	@echo "# MySQL Dump Schesuler API Documentation" > docs/api.md
	@echo "\n## Generated on: `date`\n" >> docs/api.md
	@go doc -all ./initializers/ >> docs/api.md
	@echo "Documentation generated at docs/api.md"