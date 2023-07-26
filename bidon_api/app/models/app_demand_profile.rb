class AppDemandProfile < Sequel::Model
  KEY_TO_ACCOUNT_TYPE = {
    'admob'      => 'DemandSourceAccount::Admob',
    'applovin'   => 'DemandSourceAccount::Applovin',
    'bidmachine' => 'DemandSourceAccount::BidMachine',
    'dtexchange' => 'DemandSourceAccount::DtExchange',
    'meta'       => 'DemandSourceAccount::Meta',
    'mintegral'  => 'DemandSourceAccount::Mintegral',
    'mobilefuse' => 'DemandSourceAccount::MobileFuse',
    'unityads'   => 'DemandSourceAccount::UnityAds',
    'vungle'     => 'DemandSourceAccount::Vungle',
    'bigoads'    => 'DemandSourceAccount::BigoAds',
  }.freeze

  ADAPTERS_LIST = Set.new(KEY_TO_ACCOUNT_TYPE.keys)

  many_to_one :demand_source_account, key: :account_id
end
