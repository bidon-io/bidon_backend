class MmpAccountResource < Avo::BaseResource
  self.title = :human_name

  field :id, as: :id
  field :human_name, as: :text, required: true
  field :account_type,
        as:       :select,
        required: true,
        options:  {
          Adjust:    'adjust',
          Appsflyer: 'appsflyer',
        }
  field :is_global_account, as: :boolean, required: true
  field :use_s3, as: :boolean
  field :s3_access_key_id, as: :text
  field :s3_secret_access_key, as: :text
  field :s3_bucket_name, as: :text
  field :s3_region, as: :text
  field :s3_home_folder, as: :text
  field :master_api_token, as: :text
  field :user_token, as: :text
end
