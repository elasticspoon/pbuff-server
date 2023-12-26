import ws from "k6/ws";
import { check } from "k6";

const data = `The Core Web Vitals are one of Google's Page Experience Signals
. A positive page experience naturally leads to better quality and better search engine rankings. These golden metrics help you understand which areas of your frontend application need optimization so your pages can rank higher than similar content.
Existing browser measures, such as Load and DOMContentLoaded times, no longer accurately reflect user experience very well. Relying on these load events does not give the correct metric to analyze critical performance bottlenecks that your page might have. Google's Web Vitals is a better measure of your page performance and its user experience.
`;

export const options = {
  vus: 10,
  duration: "10s",
};

export default function () {
  const url = "ws://127.0.0.1:8080";
  var params = { tags: { my_tag: "hello" } };

  var response = ws.connect(url, params, function (socket) {
    socket.on("open", function open() {
      console.log("connected");
      socket.send(data);
    });

    socket.on("message", function incoming(data) {
      socket.setTimeout(function () {
        socket.send(data);
      }, 200);
    });

    socket.on("close", function close() {
      console.log("disconnected");
    });

    socket.on("error", function (e) {
      if (e.error() != "websocket: close sent") {
        console.log("An unexpected error occurred: ", e.error());
      }
    });

    socket.setTimeout(function () {
      console.log("10 seconds passed, closing the socket");
      socket.close();
    }, 10000);
  });

  check(response, { "status is 101": (r) => r && r.status === 101 });
}
