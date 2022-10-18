class AddReplicaIdentityFullToPocTables < ActiveRecord::Migration[7.0]
  def up
    safety_assured {
      execute('ALTER TABLE app_demand_profiles REPLICA IDENTITY FULL')
      execute('ALTER TABLE app_mmp_profiles REPLICA IDENTITY FULL')
      execute('ALTER TABLE auction_configurations REPLICA IDENTITY FULL')
      execute('ALTER TABLE countries REPLICA IDENTITY FULL')
      execute('ALTER TABLE demand_sources REPLICA IDENTITY FULL')
      execute('ALTER TABLE demand_source_accounts REPLICA IDENTITY FULL')
      execute('ALTER TABLE line_items REPLICA IDENTITY FULL')
      execute('ALTER TABLE users REPLICA IDENTITY FULL')
    }
  end

  def down
    safety_assured {
      execute('ALTER TABLE app_demand_profiles REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE app_mmp_profiles REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE auction_configurations REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE countries REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE demand_sources REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE demand_source_accounts REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE line_items REPLICA IDENTITY DEFAULT')
      execute('ALTER TABLE users REPLICA IDENTITY DEFAULT')
    }
  end
end
