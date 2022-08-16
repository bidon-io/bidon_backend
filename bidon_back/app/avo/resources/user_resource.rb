class UserResource < Avo::BaseResource
  self.title = :email
  self.includes = [:apps]

  field :id, as: :id
  field :email, as: :text, required: true
end
