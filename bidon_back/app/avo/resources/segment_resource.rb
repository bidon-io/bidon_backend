class SegmentResource < Avo::BaseResource
  self.title = :name
  self.includes = [:app]

  field :id, as: :id
  field :name, as: :text
  field :description, as: :textarea
  field :filters, as: :text
  field :enabled, as: :boolean
  field :app, as: :belongs_to, required: true
end
