class DemandSourceAccount::MobileFuse < DemandSourceAccount
  def slug
    "mobile_fuse_account_#{id}"
  end
end
