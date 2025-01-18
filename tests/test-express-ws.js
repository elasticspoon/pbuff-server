import ws from "k6/ws";
import { check } from "k6";

const data = `The Core Web Vitals are one of Google's Page Experience Signals
. A positive page experience naturally leads to better quality and better search engine rankings. These golden metrics help you understand which areas of your frontend application need optimization so your pages can rank higher than similar content.
Existing browser measures, such as Load and DOMContentLoaded times, no longer accurately reflect user experience very well. Relying on these load events does not give the correct metric to analyze critical performance bottlenecks that your page might have. Google's Web Vitals is a better measure of your page performance and its user experience.
`;

function join_command(id) {
  return JSON.stringify({
    command: "join",
    message: null,
    id: id,
  });
}

function leave_command(id) {
  return JSON.stringify({
    command: "leave",
    message: null,
    id: id,
  });
}

function send_command(id) {
  return JSON.stringify({
    command: "message",
    message: data,
    id: id,
  });
}

export const options = {
  vus: 10000,
  duration: "300s",
};

const rooms = options.vus / 2;
const url = "ws://127.0.0.1:8080";
const params = { tags: { my_tag: "hello" } };

export function setup() {
  ws.connect(url, params, function (socket) {
    socket.on("open", function open() {
      for (let i = 0; i < rooms; i++) {
        socket.send(join_command(i));
      }
      socket.close();
    });
  });
}

export default function (data) {
  let room = Math.floor(Math.random() * rooms);
  var response = ws.connect(url, params, function (socket) {
    socket.on("open", function open() {
      socket.send(join_command(room));
      socket.setInterval(() => {
        socket.send(send_command(room));
      }, 1000);
    });

    socket.on("message", function incoming(data) {
      // console.log(data);
    });

    socket.on("close", function close() {
      // console.log("disconnected");
    });

    socket.on("error", function (e) {
      if (e.error() != "websocket: close sent") {
        console.log("An unexpected error occurred: ", e.error());
      }
    });

    let duration = parseInt(options.duration);
    socket.setTimeout(
      function () {
        // console.log("10 seconds passed, closing the socket");
        socket.close();
      },

      parseInt(options.duration) * 60 * 1000,
    );
  });

  check(response, { "status is 101": (r) => r && r.status === 101 });
}

export function teardown() {
  ws.connect(url, params, function (socket) {
    socket.on("open", function open() {
      for (let i = 0; i < rooms; i++) {
        leave_command(i);
      }
      socket.close();
    });
  });
}
