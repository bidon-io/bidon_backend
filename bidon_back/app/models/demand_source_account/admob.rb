class DemandSourceAccount::Admob < DemandSourceAccount
  def slug
    "admob_account_#{id}"
  end
end
