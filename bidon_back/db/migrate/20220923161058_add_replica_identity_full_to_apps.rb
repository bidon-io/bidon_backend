class AddReplicaIdentityFullToApps < ActiveRecord::Migration[7.0]
  def up
    safety_assured {
      execute('ALTER TABLE apps REPLICA IDENTITY FULL')
    }
  end

  def down
    safety_assured {
      execute('ALTER TABLE apps REPLICA IDENTITY DEFAULT')
    }
  end
end
