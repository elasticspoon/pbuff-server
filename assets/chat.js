let form = document.querySelector('#send');
let url = `ws://${location.host}${location.pathname.replace("chat", "ws")}`;
console.log(url);
let ws = new WebSocket(url);

function createMessage(message) {
  let data = JSON.parse(message);
  console.log(data);

  let container = document.createElement('div');
  container.style.display = "flex";
  container.style.gap = "10px";

  let msg = document.createElement('span');
  msg.textContent = data.body;

  let time = document.createElement('span');
  let d = new Date(data.time);
  time.textContent = d.toLocaleString() + ":";

  container.appendChild(time);
  container.appendChild(msg);

  return container;
}

function addMessage(event) {
  let message = event.data;
  let messages = document.querySelector('#chat');
  messages.appendChild(createMessage(message));
}

function sendMessage(e) {
  e.preventDefault();
  let msg = {
    body: form.elements.message.value,
    time: Date.now()
  };
  ws.send(JSON.stringify(msg));
}

ws.addEventListener("open", () => {
  console.log("Connected to the server");
  ws.addEventListener("message", addMessage);
  form.addEventListener("submit", sendMessage);
});

