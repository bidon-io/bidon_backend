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

  validates :bid_floor, numericality: { greater_than_or_equal_to: 0 }
  validates :format,
            presence: { message: 'must be present for banner' }, # rubocop:disable Rails/I18nLocaleTexts
            if:       :banner?
  validates :format,
            absence: { message: ->(item, _) { "must be blank for #{item.ad_type}" } },
            unless:  :banner?

  def extra=(value)
    if value.is_a?(Hash)
      super(value)
    else
      super(JSON.parse(value.gsub('=>', ':')))
    end
  end
end
