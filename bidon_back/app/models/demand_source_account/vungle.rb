class DemandSourceAccount::Vungle < DemandSourceAccount
  def slug
    "vungle_account_#{id}"
  end
end
