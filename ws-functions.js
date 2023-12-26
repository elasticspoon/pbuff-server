module.exports = { createSubCommand, createSendCommand };

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

function createSubCommand(userContext, events, done) {
  userContext.vars.sub_cmd = JSON.stringify(sub_cmd);
  return done();
}

function createSendCommand(userContext, events, done) {
  userContext.vars.send_cmd = JSON.stringify(send_command);
  return done();
}
