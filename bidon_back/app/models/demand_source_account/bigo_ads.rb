class DemandSourceAccount::BigoAds < DemandSourceAccount
  def slug
    "bigoads_account_#{id}"
  end
end
