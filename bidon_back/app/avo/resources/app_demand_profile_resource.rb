class AppDemandProfileResource < Avo::BaseResource
  self.title = :slug

  field :id, as: :id
  field :app, as: :belongs_to, required: true
  field :demand_source, as: :belongs_to, required: true
  field :account, as: :belongs_to, required: true
  field :data, as: :code
  field :account_type, as: :select, required: true, options: ::DemandSourceType::OPTIONS
end
