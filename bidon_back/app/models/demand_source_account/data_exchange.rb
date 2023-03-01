class DemandSourceAccount::DataExchange < DemandSourceAccount
  def slug
    "data_exchange_account_#{id}"
  end
end
