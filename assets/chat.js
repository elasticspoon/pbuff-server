let form = document.querySelector('#send');
let url = `ws://${location.host}${location.pathname.replace("chat", "ws")}`;
console.log(url);
let ws = new WebSocket(url);

function addMessage(event) {
  let message = event.data;
  let messages = document.querySelector('#chat');
  let span = document.createElement('span');
  span.style.display = "block";
  span.textContent = message;
  messages.appendChild(span);
}

function sendMessage(e) {
  e.preventDefault();
  ws.send(form.elements.message.value);
}

ws.addEventListener("open", () => {
  console.log("Connected to the server");
  ws.addEventListener("message", addMessage);
  form.addEventListener("submit", sendMessage);
});

