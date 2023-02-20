class CountryResource < Avo::BaseResource
  self.title = :human_name

  field :id, as: :id
  field :human_name, as: :text, required: true
  field :alpha_2_code, as: :text, required: true
  field :alpha_3_code, as: :text, required: true
end
