# websocket-api-gateway
Websocket API Gateway that allows to subscribe on notifications about price changes of financial instruments

To test connection with server, send such request:

/ping -method GET, it should return such JSON:

{

    "message": "pong"

}

To connect to websocket for receiving notifications from Bitmex, send such request:

/ws -method GET

To subscribe on notifications from Bitmex, send such JSON message to websocket:

{

    "action": "subscribe", // required

    "symbols": <[]string>  // optional, list of trade instruments, in case of absence subscription will be on all instruments

}

To unsubscribe from notifications from Bitmex, send such JSON message to websocket:

{

    "action": "unsubscribe" // required

}

