# chatRPC
**Description**: backend of chat application, which makes it possible to connect to any room and communicate 
directly with other users.
If you want to improve rpc methods, use this command for codegen:
``` bash 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protofiles/chat.proto
```

### 1.Download and run the server
#### Clone repository:
``` bash 
git clone https://github.com/sergripenko/chatRPC.git
```

#### Install dependencies:
``` bash 
go mod vendor
```

#### Open project repo and run server:
``` bash 
go run cmd/main.go
```
All logs will be saved in **logs.txt** 

### 2.Run clients
#### Connect to server:
``` bash 
go run client/connect/main.go -username {username}
```

#### Join group chat:
``` bash 
go run client/join_group/main.go -username {username} -group {group name}
```

#### Leave group chat:
``` bash 
go run client/leave_group/main.go -username {username} -group {group name}
```

#### Create group chat:
``` bash 
go run client/create_group/main.go -username {username} -group {group name}
```

#### Send direct message to user:
``` bash 
go run client/send_message/main.go -username {from username} -user {to username} -message {message text}
```

#### Send message to group:
``` bash 
go run client/send_message/main.go -username {from username} -group {group name} -message {message text}
```

#### List channels:
``` bash 
go run client/list_channels/main.go -username {username}
```


