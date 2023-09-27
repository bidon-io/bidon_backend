class RemoveNotNullContraintFromUserIdDemandSourceAccounts < ActiveRecord::Migration[7.0]
  def change
    change_column_null :demand_source_accounts, :user_id, true
  end
end
