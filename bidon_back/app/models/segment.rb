class Segment < ApplicationRecord
  belongs_to :app

  def filters=(value)
    if value.is_a?(Array)
      super(value)
    else
      super(JSON.parse(value.gsub('=>', ':')))
    end
  end
end
