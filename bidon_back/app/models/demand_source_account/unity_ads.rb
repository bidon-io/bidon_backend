class DemandSourceAccount::UnityAds < DemandSourceAccount
  def slug
    "unity_ads_account_#{id}"
  end
end
