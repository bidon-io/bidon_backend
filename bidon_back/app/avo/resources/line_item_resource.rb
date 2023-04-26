class LineItemResource < Avo::BaseResource
  self.title = :human_name
  self.keep_filters_panel_open = true

  field :id, as: :id
  field :human_name, as: :text, required: true
  field :app, as: :belongs_to, required: true
  field :bid_floor, as: :number, required: true
  field :ad_type, as: :select, required: true, enum: ::AuctionConfiguration.ad_types
  field :format,
        as:            :select,
        enum:          ::LineItem.formats,
        display_value: true,
        include_blank: 'N/A'
  field :account, as: :belongs_to, required: true
  field :account_type, as: :select, required: true, options: ::DemandSourceType::OPTIONS
  field :code, as: :text
  field :extra, as: :code, hide_on: :index

  filter AdTypeFilter
  filter AppFilter
end
