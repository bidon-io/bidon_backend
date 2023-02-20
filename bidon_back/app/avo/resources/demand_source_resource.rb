class DemandSourceResource < Avo::BaseResource
  self.title = :human_name

  field :id, as: :id
  field :human_name, as: :text, required: true
  field :api_key, as: :text, required: true
end
