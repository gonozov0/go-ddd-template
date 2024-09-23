# go-echo-template

The project is a template for Go repositories using the Echo framework.
Contains examples of the code structure for a project with a domain-driven design (DDD) approach.


## Настойка dev окружения

### Go

goenv — это инструмент для управления версиями Go, который поддерживает множество операционных систем, включая Windows.

```sh
brew install goenv
```

После установки вам нужно будет инициализировать `goenv`. Добавьте следующие строки в ваш файл профиля (`~/.bash_profile`, `~/.zshrc`, `~/.profile` и т.д.), чтобы инициализировать `goenv` при каждом открытии терминала:

```sh
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
```

Не забудьте применить изменения в файле профиля, выполнив команду `source ~/.bash_profile` (или аналогичную, в зависимости от вашей оболочки).

После установки `goenv`, вы можете устанавливать разные версии Go и переключаться между ними. Для установки новой версии используйте:
```sh
goenv install 1.x.x
```
А для выбора версии Go для текущего проекта или глобально для пользователя:
```sh
goenv local 1.x.x
goenv global 1.x.x
```

Чтобы настроить конкретную версию `Go` для проекта в IDE (например GoLand), выберите в качестве `GOROOT` нужную версию в директории `~/.goenv/versions`

### Линтеры и кодогенерация

Для запуска `make lint` и `go generate ./...` необходимо установить следующие утилиты:

```sh
brew install protobuf
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/segmentio/golines@latest
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

либо просто выполнить команду `make install`

### Pre-commit hooks

Для установки pre-commit hooks выполните команды:

```sh
brew install pre-commit
pre-commit install
```

### Golang-migrate

Для выполнения миграций базы данных используется утилита `golang-migrate`.
Чтобы установить утилиту, выполните команду:

```shell
brew install golang-migrate
```

Чтобы применить миграции к базе, существует make команда:

```shell
make migrate_up
```

Также для отката миграций (параметр `count` указывает количество миграций, которые нужно откатить):
```shell
make migrate_down count=1
```

Для создания новой миграции используйте команду:
```shell
make create_migration name=migration_name
```
