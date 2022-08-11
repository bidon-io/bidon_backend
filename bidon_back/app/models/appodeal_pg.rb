# SHOULD NOT BE USED FOR BUSINESS LOGIC
# ADDED JUST FOR POC, SO WE CAN SYNC DATA FROM APPODEAL FASTER
# WE SHOULD USE KAFKA IN FUTURE AND REMOVE HARD DEPENDENCY ON APPODEAL
class AppodealPg < ApplicationRecord
  self.abstract_class = true

  connects_to database: { reading: :appodeal_read_only }

  def self.execute(sql)
    connected_to(role: :reading, prevent_writes: true) { connection.execute(sql).entries }
  end
end
