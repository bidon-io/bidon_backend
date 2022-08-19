class AuctionConfigurationResource < Avo::BaseResource
  self.title = :name
  self.includes = [:app]

  field :id, as: :id
  field :app, as: :belongs_to, required: true
  field :name, as: :text, required: true
  field :ad_type, as: :select, required: true, enum: ::AuctionConfiguration.ad_types
  field :rounds, as: :textarea, required: true
end
