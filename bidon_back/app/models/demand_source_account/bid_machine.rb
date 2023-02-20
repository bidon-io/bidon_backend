class DemandSourceAccount::BidMachine < DemandSourceAccount
  def slug
    "bidmachine_account_#{id}"
  end
end
