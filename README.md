# App
## List of contents
- [Logger](#logger)
    - [yaml](#yaml)
    - [env](#env)
    - [args](#args)

## Logger
### yaml
```yaml
logger:
  level: debug
  writeToFile: true
  file:
    path: ./
    name: app.log
    maxAge: 24h
    rotationTime: 168h
```
### env
```bash
LOGGER_LEVEL="info" # [LOG_LEVEL, LOGGING]
LOGGER_WRITE_TO_FILE=false # [LOG_WRITE_TO_FILE]
LOGGER_FILE_PATH="~/" # [LOG_FILE_PATH]
LOGGER_FILE_NAME="app.log" # [LOG_FILE_NAME]
LOGGER_FILE_MAX_AGE="24h" # [LOG_FILE_MAX_AGE]
LOGGER_FILE_ROTATION_TIME="168h" # [LOG_FILE_ROTATION_TIME]
```
### args
```bash
$ ./build/app --logger.level info
```