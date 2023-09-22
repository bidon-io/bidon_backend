class DemandSourceAccount::Amazon < DemandSourceAccount
  def slug
    "amazon_account_#{id}"
  end
end
