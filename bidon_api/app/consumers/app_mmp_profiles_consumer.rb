# frozen_string_literal: true

# Example consumer that prints messages payloads
class AppMmpProfilesConsumer < ApplicationConsumer
  def params_for(data)
    params = super(data)
    params['start_date'] = (Time.zone.at(0).to_date + params['start_date']).to_date
    params
  end

  def model
    AppMmpProfile
  end
end
