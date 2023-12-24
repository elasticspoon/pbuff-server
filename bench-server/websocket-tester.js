const WebSocket = require("ws");

let ws = new WebSocket("ws://localhost:3001/cable");
let id = { channel: "ChatRoomChannel" };
let sub_cmd = {
  command: "subscribe",
  identifier: JSON.stringify(id),
};

let send_command = {
  command: "message",
  identifier: JSON.stringify(id),
  data: JSON.stringify({ content: "Hello World" }),
};

let num = 0;

ws.onclose = function () {
  // thing to do on close
};

ws.onerror = function () {
  // thing to do on error
};

ws.onmessage = function (e) {
  // thing to do on message
  num += 1;
  ws.send(JSON.stringify(send_command));

  if (num > 10) {
    ws.close();
  }
};

ws.onopen = function (e) {
  ws.send(JSON.stringify(sub_cmd));
};
