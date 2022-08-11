class PostResource < Avo::BaseResource
  self.model_class = ActiveRecord::SchemaMigration
end
