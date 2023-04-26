class AddFormatToBanners < ActiveRecord::Migration[7.0]
  LineItem = Class.new(ActiveRecord::Base) # rubocop:disable Rails/ApplicationRecord

  FORMATS = {
    '320x50'  => 'BANNER',
    '728x90'  => 'LEADERBOARD',
    '300x250' => 'MREC',
    '0x50'    => 'ADAPTIVE',
  }.freeze

  def up
    add_column :line_items, :format, :string

    LineItem.where(ad_type: 3).find_each do |line_item|
      format = FORMATS["#{line_item.width}x#{line_item.height}"]
      next unless format

      line_item.update!(format:)
    end
  end

  def down
    remove_column :line_items, :format
  end
end
