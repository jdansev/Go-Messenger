


// fuzzy find hubs

ws = new WebSocket("ws://localhost:1212/find-hubs/123");

ws.send("hello");

ws.onmessage = function (evt) {
    var messages = evt.data.split('\n');
    var msg = JSON.parse(messages[0]);
    console.log(msg);
};