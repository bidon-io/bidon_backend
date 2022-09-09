class AppResource < Avo::BaseResource
  self.title = :package_name

  field :id, as: :id
  field :platform_id, as: :number, required: true
  field :human_name, as: :text, required: true
  field :package_name, as: :text
  field :app_key, as: :text
  field :settings, as: :code
end
