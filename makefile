TARGET_DIR=public

all: setup clean
	go run .

setup:
	mkdir -p $(TARGET_DIR)

clean:
	rm -rf $(TARGET_DIR)/*

watch: all
	fswatch --latency 0.1 --one-per-batch -0 . --exclude=$(TARGET_DIR) | xargs -0 -n1 -I{} make

serve:
	browser-sync --no-ui --no-notify --server --watch $(TARGET_DIR) $(TARGET_DIR)

watch-serve: all
	# Spawn in parallel browser-sync to serve the files and fswatch to re-make when
	# files change. Kill both with one <CTRL+C>.
	(trap 'kill 0' SIGINT; browser-sync --no-ui --no-notify --server --watch $(TARGET_DIR) $(TARGET_DIR) & fswatch --latency 0.1 --one-per-batch -0 . --exclude=$(TARGET_DIR) | xargs -0 -n1 -I{} make)
