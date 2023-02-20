class DemandSourceAccount::Applovin < DemandSourceAccount
  def slug
    "applovin_account_#{id}"
  end
end
