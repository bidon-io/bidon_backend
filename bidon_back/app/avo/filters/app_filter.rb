class AppFilter < Avo::Filters::SelectFilter
  self.name = "App filter"

  def apply(_request, query, value)
    query = query.where(app_id: value) if value.present?
    query
  end

  def options
    App.select(:id, :package_name).each_with_object({}) do |app, options|
      options[app.id] = app.package_name
    end
  end
end
