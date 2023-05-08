class DemandSourceAccount::DtExchange < DemandSourceAccount
  def slug
    "dtexchange_account_#{id}"
  end
end
