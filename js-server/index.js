"use strict";

const express = require("express");
const path = require("path");
const { createServer } = require("http");

const { WebSocketServer } = require("ws");

const app = express();
app.use(express.static(path.join(__dirname, "/public")));

const server = createServer(app);
const wss = new WebSocketServer({ server });

wss.on("connection", function (ws) {
  console.log("Client connected");
  ws.on("message", function message(data, isBinary) {
    wss.clients.forEach(function each(client) {
      if (client.readyState === 1) {
        client.send(data, { binary: isBinary });
      }
    });
  });

  ws.on("error", console.error);

  ws.on("close", function () {
    console.log("closing connection");
  });
});

server.listen(8080, function () {
  console.log("Listening on http://localhost:8080");
});
