class DemandSourceAccount::Inmobi < DemandSourceAccount
  def slug
    "inmobi_account_#{id}"
  end
end
