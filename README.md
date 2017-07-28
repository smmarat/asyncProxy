# AsyncProxy

## Build

go build proxy.go

## Usage

./proxy -url=callback -port=port

callback - url to receive http callback from async service
<br/>port - port to listen (9090 by default) 

### Example

- You want to interrupt with your process in corezoid.com 
- Your process webhook is https://www.corezoid.com/api/1/json/public/255213/e0ac428061effa899d78f4085fe3983f0b96dfd5
- You don't have static ip in your application (mobile deice, etc.)
- You started proxy binary in your vps and specify callback url
- You sent json POST request to your vps and see proxy logs: 

got a request!
<br/>Request: {
<br/>  "webhook":"https://www.corezoid.com/api/1/json/public/255213/e0ac428061effa899d78f4085fe3983f0b96dfd5",
<br/>  "data":{"a":"aa","b":"bb"}
<br/>}
<br/>Sync response: {"request_proc":"ok","ops":{"proc":"ok","obj":"task","ref":null,"obj_id":"597ad4d8f6c3767e8da7be19"}}
<br/>got a callback!
<br/>Request: {"ref":"337188384","x":"YYY","sys":{"ref":null,"obj_id":"597ad4d8f6c3767e8da7be19","conv_id":255213,"node_id":"5979c38b60e32776634c619a"}}
<br/>Async response: {"ref":"337188384","x":"YYY","sys":{"ref":null,"obj_id":"597ad4d8f6c3767e8da7be19","conv_id":255213,"node_id":"5979c38b60e32776634c619a"}}

- Proxy will wait async response from Corezoid and will send it to your application as sync response