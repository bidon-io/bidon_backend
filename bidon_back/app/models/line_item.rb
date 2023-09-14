class LineItem < ApplicationRecord
  self.ignored_columns += %w[width height]

  belongs_to :app
  belongs_to :account, class_name: 'DemandSourceAccount'

  enum ad_type: AdType::ENUM
  enum format: {
    banner:      'BANNER',
    leaderboard: 'LEADERBOARD',
    mrec:        'MREC',
    adaptive:    'ADAPTIVE',
  }, _prefix: true

  def extra=(value)
    if value.is_a?(Hash)
      super(value)
    else
      super(JSON.parse(value.gsub('=>', ':')))
    end
  end
end
