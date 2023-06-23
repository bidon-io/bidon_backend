class AuctionConfigurationResource < Avo::BaseResource
  self.title = :name
  self.includes = [:app]

  field :id, as: :id
  field :app, as: :belongs_to, required: true
  field :name, as: :text, required: true
  field :ad_type, as: :select, required: true, enum: ::AuctionConfiguration.ad_types
  field :pricefloor, as: :number, required: true
  field :rounds, as: :code, required: true
  field :external_win_notifications, as: :boolean, required: true
end
