# AsyncProxy

## Build

go build proxy.go

## Usage

./proxy -url=callback -port=port

callback - url to receive http callback from async service
<br/>port - port to listen (9090 by default) 

### Example

- You want to interact with your process in corezoid.com 
- Your process webhook is https://www.corezoid.com/api/1/json/public/255213/e0ac428061effa899d78f4085fe3983f0b96dfd5
- You don't have static ip in your application (mobile device, etc.)
- You started proxy binary in your vps and specify callback url
- You sent json POST request to your vps and see proxy logs: 

```
got a request!
Request: {
  "webhook":"https://www.corezoid.com/api/1/json/public/255213/e0ac428061effa899d78f4085fe3983f0b96dfd5",
  "data":{"a":"aa","b":"bb"}
}
Sync response: {"request_proc":"ok","ops":{"proc":"ok","obj":"task","ref":null,"obj_id":"597ad4d8f6c3767e8da7be19"}}
got a callback!
Request: {"ref":"337188384","x":"YYY","sys":{"ref":null,"obj_id":"597ad4d8f6c3767e8da7be19","conv_id":255213,"node_id":"5979c38b60e32776634c619a"}}
Async response: {"ref":"337188384","x":"YYY","sys":{"ref":null,"obj_id":"597ad4d8f6c3767e8da7be19","conv_id":255213,"node_id":"5979c38b60e32776634c619a"}}
```
Your request should contain your data (with be used in Corezoid logic) in field "data" and your process webhook in field "webhook"   

- Proxy will wait asynchronous response from Corezoid and will send it to your application as synchronous response