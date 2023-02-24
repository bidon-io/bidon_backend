class LineItemResource < Avo::BaseResource
  self.title = :human_name

  field :id, as: :id
  field :human_name, as: :text, required: true
  field :app, as: :belongs_to, required: true
  field :bid_floor, as: :number, required: true
  field :ad_type, as: :select, required: true, enum: ::AuctionConfiguration.ad_types
  field :account, as: :belongs_to, required: true
  field :account_type, as: :select, required: true, options: ::DemandSourceType::OPTIONS
  field :code, as: :text
  field :height, as: :number, hide_on: :index
  field :width, as: :number, hide_on: :index
  field :extra, as: :code, hide_on: :index
end
