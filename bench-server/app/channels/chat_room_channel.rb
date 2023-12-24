class ChatRoomChannel < ApplicationCable::Channel
  def subscribed
    stream_from "test"
  end

  def unsubscribed
    # Any cleanup needed when channel is unsubscribed
  end

  def receive(data)
    ActionCable.server.broadcast("test", data)
  end
end
