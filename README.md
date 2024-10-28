# Web-socket Comment Board Service
This service provides broadcasting messages to clients subscribed to a board and as well as sending one-to-one messages to clients 
### Useful links
`WebSocket Endpoint` - [Subscribe to board](prayer-service-495160257238.us-east4.run.app/subscribe)

`Swagger Docs` - [Swagger Docs](https://prayer-service-495160257238.us-east4.run.app/docs/index.html)

`Postman workspace` - [Join Workspace to view test routes](https://app.getpostman.com/join-team?invite_code=848185c06cac47bde6108f501efd035e) 

`Go on Azure` - [Hosting Go on Azure](https://azure.microsoft.com/en-us/resources/developers/go-apps)

### Headers

| Name    | Type | Required    | Description |
| -------- | ------- | -------- | ------- |
| Token  | string    | Yes  | User Auth Token    |
| Board | string     | Yes  | Board Id to connect to |


Specifies the type of device making the request. Defaults to "web" if not provided.



Available events types and actions: 

```go
const (
	// Event types 
	PRAYER = "PRAYER"
	COMMENT = "COMMENT"
	BOARD = "BOARD"
	// Event actions
	ADD = "ADD"
	DELETE = "DELETE"
	UPDATE = "UPDATE"
	SWITCH = "SWITCH"
)
```

When connect clients, they can send event payload: 

```go
type Event struct {
	Type string `json:"type"`
	Action string `json:"action"`
	Data string `json:"data"`
}
```
i.e an example payload of adding a prayer:

```json
{
    "type": "PRAYER",
    "action": "ADD",
    "data": "{\n        \"id\": 2,\n        \"boardId\": 1,\n        \"title\": \"Healing for my dog\",\n        \"description\": \"I need prayer for my dog, he is sick\",\n        \"comments\": [\n            {\n                \"comment\": \"He will be healed in Jesus name\",\n                \"id\": 1,\n                \"prayer_id\": 2,\n                \"user\": {\n                    \"userId\": \"1\",\n                    \"firstName\": \"John\",\n                    \"lastName\": \"Doe\",\n                    \"userName\": \"johndoe\",\n                    \"screenName\": \"John Doe\",\n                    \"email\": \"\"\n                }\n            }\n        ]\n    }"
}
```

i.e same payload can be used delete a prayer `NOTE` the action would change to `DELETE`:

```json
{
    "type": "PRAYER",
    "action": "DELETE",
    "data": "{\n        \"id\": 2,\n        \"boardId\": 1,\n        \"title\": \"Healing for my dog\",\n        \"description\": \"I need prayer for my dog, he is sick\",\n        \"comments\": [\n            {\n                \"comment\": \"He will be healed in Jesus name\",\n                \"id\": 1,\n                \"prayer_id\": 2,\n                \"user\": {\n                    \"userId\": \"1\",\n                    \"firstName\": \"John\",\n                    \"lastName\": \"Doe\",\n                    \"userName\": \"johndoe\",\n                    \"screenName\": \"John Doe\",\n                    \"email\": \"\"\n                }\n            }\n        ]\n    }"
}
```
