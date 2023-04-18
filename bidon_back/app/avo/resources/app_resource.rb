class AppResource < Avo::BaseResource
  self.title = :slug

  field :id, as: :id
  field :platform_id, as: :select, required: true, enum: ::App.platform_ids
  field :human_name, as: :text, required: true
  field :package_name, as: :text
  field :user, as: :belongs_to, required: true
  field :app_key, as: :text
  field :settings, as: :code

  field :line_items, as: :has_many
end
