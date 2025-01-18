"use strict";

const express = require("express");
const path = require("path");
const { createServer } = require("http");

const { WebSocketServer } = require("ws");

const app = express();
app.use(express.static(path.join(__dirname, "/public")));

const server = createServer(app);
const wss = new WebSocketServer({ server });

const chatRooms = new Map();
let chatRoomId = 1;

wss.on("connection", function (ws) {
  console.log("Client connected");

  function leaveChatRoom(id) {
    let cr = chatRooms.has(id) ? chatRooms.get(id) : null;
    if (cr) {
      cr.indexOf(ws);
      cr.splice(cr.indexOf(ws), 1);
    }
  }

  ws.on("message", function message(data, isBinary) {
    let parsedData = JSON.parse(data);
    const { message, command, id } = parsedData;

    if (command === "join") {
      let cr = chatRooms.get(id);
      chatRooms.has(id) ? chatRooms.get(id).push(ws) : chatRooms.set(id, [ws]);
    } else if (command === "leave") {
      // console.log("Leaving room", room);
      leaveChatRoom(id);
    } else {
      let cr = chatRooms.has(id) ? chatRooms.get(id) : null;
      // console.log(cr, chatRooms, id, parsedData);
      if (cr) {
        cr.forEach(function each(client) {
          client.send(data, { binary: isBinary });
        });
      }
    }
  });

  ws.on("error", console.error);

  ws.on("close", function () {
    console.log("Client disconnected");
  });
});

server.listen(8080, function () {
  console.log("Listening on http://localhost:8080");
});
