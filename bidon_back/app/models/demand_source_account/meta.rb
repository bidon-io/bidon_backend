class DemandSourceAccount::Meta < DemandSourceAccount
  def slug
    "meta_account_#{id}"
  end
end
