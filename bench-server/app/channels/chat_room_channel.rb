class ChatRoomChannel < ApplicationCable::Channel
  def subscribed
    stream_from "test"
  end

  def unsubscribed
    puts "unsubscribed"
  end

  def receive(data)
    ActionCable.server.broadcast("test", data)
  end
end
