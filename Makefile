run:
	echo "start builind..." && \
	docker compose -f docker-compose.yml build && \
	echo "run docker compose..." && \
	docker compose up -d && \
	echo "DONE!"

stop:
	echo "stop project" && \
	docker compose -f docker-compose.yml down
	echo "DONE!"