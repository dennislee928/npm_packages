# BitYacht Exchange Api Server

## 資料夾架構

```
.
├── README.md              # 說明文件
├── api
│   ├── controllers
│   │   └── v1
│   ├── middelware
│   └── router.go
├── cmd
├── configs
│   ├── config.yaml         # 設定檔
│   ├── config.yaml.example # 設定檔範例
│   └── configs.go          # Viper 設定檔
├── docs
│   ├── docs.go            # Swagger 設定檔
│   ├── swagger.json       # Swagger JSON
│   └── swagger.yaml       # Swagger YAML
└── internal
    ├── cache               # redis
    │   └── redis
    ├── database            # gorm
    │   └── sql
    └── pkg                 # 共用套件
        └── err             # 自定義錯誤
```

## 安裝注意事項

1. Recommend Versions
   - Golang: 1.20.5
   - MariaDB: 10.11.4
   - Redis: 7.0.12
1. Set Time Zone Path for MariaDB
   ```bash
    mysql_tzinfo_to_sql /usr/share/zoneinfo | mysql -u root mysql -p
   ```
1. Build the executable file
   ```bash
   go build
   ```
1. Environment Variable

   Copy `config.example.yaml` and Rename it to `config.yaml`. Then you can Modify it by Real Environment.

   Notice: the `config.yaml` can place in same path of executable file or `./configs/config.yaml`

   Example:

   ```
    Usually:
    /.../bityacht-exchange-api-server # executable file
    /.../config.yaml

    or

    /.../bityacht-exchange-api-server # executable file
    /.../configs/config.yaml
   ```

1. Database Migration & Seeder

   Assum the name of executable file is `bityacht-exchange-api-server`.

   ```bash
   ./bityacht-exchange-api-server migrate up --create
   ./bityacht-exchange-api-server seed default
   ```

   `migrate up --create` will check the database exist or not, it will create the database if not exist.

   `seed default` will create the default records to tables.

1. Check the Permission of Storage Path

   The default Storage Path is `/var/bityacht_exchange` (can be modify in `config.yaml`)

   Make sure the Process have Permission to Access (Read & Write) this path.
