import ws from "k6/ws";
import http from "k6/http";
import { check } from "k6";

const data = `The Core Web Vitals are one of Google's Page Experience Signals
. A positive page experience naturally leads to better quality and better search engine rankings. These golden metrics help you understand which areas of your frontend application need optimization so your pages can rank higher than similar content.
Existing browser measures, such as Load and DOMContentLoaded times, no longer accurately reflect user experience very well. Relying on these load events does not give the correct metric to analyze critical performance bottlenecks that your page might have. Google's Web Vitals is a better measure of your page performance and its user experience.
`;

const message = JSON.stringify({
  body: data,
  time: 234,
});

export const options = {
  vus: 10000,
  duration: "30s",
};

let rooms = options.vus / 2;
const url = "http://localhost:3000";
const params = { tags: { my_tag: "hello" } };

export function setup() {
  for (let i = 0; i < rooms; i++) {
    const res = http.get(url + "/new");
    check(res, { "status is 200": (r) => r && r.status === 200 });
  }
}

export default function () {
  const wsUrl = "ws://localhost:3000/ws/" + Math.floor(Math.random() * rooms);

  var response = ws.connect(wsUrl, params, function (socket) {
    socket.on("open", function open() {
      // console.log("connected");
      socket.setInterval(() => {
        socket.send(message);
      }, 1000);
      // socket.send(message);
    });

    socket.on("message", function incoming(e) {
      // socket.setTimeout(function () {
      //   socket.send(message);
      // }, 200);
    });

    socket.on("close", function close() {
      // console.log("disconnected");
    });

    socket.on("error", function (e) {
      // if (e.error() != "websocket: close sent") {
      //   console.log("An unexpected error occurred: ", e.error());
      // }
    });

    socket.setTimeout(
      function () {
        socket.close();
      },
      // parseInt(options.duration) * 60 * 1000,
      parseInt(options.duration) * 1000,
    );
  });

  check(response, { "status is 101": (r) => r && r.status === 101 });
}

export function teardown() {
  for (let i = 0; i < rooms; i++) {
    const res = http.del(url + "/chat/" + i);
    check(res, { "status is 200": (r) => r && r.status === 200 });
  }
}
