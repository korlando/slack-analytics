# Slack Analytics

## Usage

Export your slack data and unzip it into a folder, say `data`.

```
go run ./cmd/analyze.go -p ./data
```

Use `-p` or `--path` to specify the path to the slack data folder.
Use `-m` to specify that the path points to a JSON file containing a messages array.
