# Use this class to generate public_uid for existing records: segments, line_items, auction_configurations, etc
class PublicUidGenerator
  def self.generate_for(record)
    return if record.public_uid.present?

    timestamp = record.created_at.strftime('%s%L').to_i
    worker_id = 1
    snowflake = TwitterSnowflake.synthesize(timestamp:, worker_id:)

    loop do
      public_uid = snowflake.id

      if record.class.where(public_uid:).count == 0
        record.update!(public_uid:)
        break
      else
        timestamp += 1
      end
    end
  end
end

# Examples:
# Segment.find_each { |r| PublicUidGenerator.generate_for(r) }
# LineItem.find_each { |r| PublicUidGenerator.generate_for(r) }
# AuctionConfiguration.find_each { |r| PublicUidGenerator.generate_for(r) }
# App.find_each { |r| PublicUidGenerator.generate_for(r) }
# DemandSourceAccount.find_each { |r| PublicUidGenerator.generate_for(r) }
# AppDemandProfile.find_each { |r| PublicUidGenerator.generate_for(r) }
# User.find_each { |r| PublicUidGenerator.generate_for(r) }
