# webhook-go

支持GitHub、Gitea、Gogs


## Configuration

> URL：http://127.0.0.1:8000/webHooks

> 复制`data`目录中的`config.example.json`文件重命名为`config.json`并修改里面的值

```json
{
  "这里是仓库full_name的值": {
    "logfile": "test-gitea-webhook.log",
    "secret": "在Webhooks中设定的secret",
    "commands": [
      "data/update_repo.sh"
    ]
  }
}
```

## 运行

```bash
# -h为地址（默认0.0.0.0），-p为端口（默认8000）
./webhook-go -h 127.0.0.1 -p 8000
```