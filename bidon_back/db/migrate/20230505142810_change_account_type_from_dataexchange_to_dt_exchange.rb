class ChangeAccountTypeFromDataexchangeToDtExchange < ActiveRecord::Migration[7.0]
  def up
    safety_assured do
      execute(
        <<~SQL.squish,
          UPDATE demand_source_accounts
          SET type = 'DemandSourceAccount::DtExchange'
          WHERE type = 'DemandSourceAccount::DataExchange'
        SQL
      )
      execute(
        <<~SQL.squish,
          UPDATE app_demand_profiles
          SET account_type = 'DemandSourceAccount::DtExchange'
          WHERE account_type = 'DemandSourceAccount::DataExchange'
        SQL
      )
    end
  end

  def down; end
end
