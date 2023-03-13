class AppFilter < Avo::Filters::SelectFilter
  self.name = 'App filter'

  def apply(_request, query, value)
    query = query.where(app_id: value) if value.present?
    query
  end

  def options
    App.select(:id, :package_name, :platform_id).each_with_object({}) do |app, options|
      options[app.id] = app.slug
    end
  end
end
