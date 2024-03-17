TARGET_DIR=public
CACHE_DIR=cache

WATCH_CMD=fswatch --latency 0.1 --one-per-batch -0 . --exclude=$(TARGET_DIR) --exclude=$(CACHE_DIR) | xargs -0 -n1 -I{} make
SERVE_CMD=browser-sync --no-ui --no-notify --server --watch $(TARGET_DIR) $(TARGET_DIR)

.PHONY: all
all: setup clean
	go run .

.PHONY: setup
setup:
	mkdir -p $(TARGET_DIR)

.PHONY: clean
clean:
	rm -rf $(TARGET_DIR)/*

.PHONY: format
format:
	gofmt -s -w .

.PHONY: watch
watch: all
	$(WATCH_CMD)

.PHONY: serve
serve:
	$(SERVE_CMD)

.PHONY: watch-serve
watch-serve: all
	# Spawn in parallel browser-sync to serve the files and fswatch to re-make when
	# files change. Kill both with one <CTRL+C>.
	(trap 'kill 0' SIGINT; $(SERVE_CMD) & $(FSWATCH_CMD))
