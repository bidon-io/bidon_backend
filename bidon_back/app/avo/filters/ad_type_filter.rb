class AdTypeFilter < Avo::Filters::SelectFilter
  self.name = 'Ad type filter'

  def apply(_request, query, value)
    query = query.where(ad_type: value) if value.present?
    query
  end

  def options
    AdType::ENUM.invert
  end
end
