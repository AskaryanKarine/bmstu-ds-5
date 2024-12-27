MIGRATION_DIR = "migrations"

create-migration:
ifeq ($(name),)
	@echo "Input service name with name param"
else ifeq ($(service),)
	@echo "Input service name with service param"
else
	goose --dir=$(MIGRATION_DIR)/$(service) create $(name) sql
endif


run-migration:
	for name in $(MIGRATION_DIR)/* ; do \
		service_name=$$(echo $${name} | sed 's|^migrations/||'); \
        service_dir=$${service_name}_service; \
        echo "Service name: $$service_dir"; \
		goose --dir=$${name} postgres 'postgresql://postgres:postgres@localhost:5432/'$$service_dir'?sslmode=disable' up; \
	done

lint:
	golangci-lint run --timeout=5m

.PHONY : create-migration lint
