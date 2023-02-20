class AppMmpProfileResource < Avo::BaseResource
  self.title = :human_name

  field :id, as: :id
  field :app, as: :belongs_to, required: true
  field :start_date, as: :date, required: true
  field :mmp_platform, as: :select, required: true, enum: AppMmpProfile.mmp_platforms
  field :mmp_account_primary, as: :belongs_to, required: true
  field :mmp_account_secondary, as: :belongs_to, required: true
  field :get_spend_from_secondary_mmp_account, as: :boolean
  field :primary_mmp_raw_data_source,
        as:       :select,
        required: true,
        enum:     AppMmpProfile.primary_mmp_raw_data_sources
  field :secondary_mmp_raw_data_source,
        as:       :select,
        required: true,
        enum:     AppMmpProfile.secondary_mmp_raw_data_sources
  field :adjust_app_token, as: :text
  field :appsflyer_dev_key, as: :text
  field :appsflyer_conversion_keys, as: :text
  field :firebase_config_keys, as: :text
  field :firebase_expiration_duration, as: :number
  field :firebase_tracking, as: :boolean
  field :facebook_tracking, as: :boolean
end
