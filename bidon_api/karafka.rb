# frozen_string_literal: true

class KarafkaApp < Karafka::App
  setup do |config|
    config.kafka = { 'bootstrap.servers': 'kafka:9092' }
    config.client_id = 'example_app'
    # Recreate consumers with each batch. This will allow Rails code reload to work in the
    # development mode. Otherwise Karafka process would not be aware of code changes
    config.consumer_persistence = !Rails.env.development?
  end

  # Comment out this part if you are not using instrumentation and/or you are not
  # interested in logging events for certain environments. Since instrumentation
  # notifications add extra boilerplate, if you want to achieve max performance,
  # listen to only what you really need for given environment.
  Karafka.monitor.subscribe(Karafka::Instrumentation::LoggerListener.new)
  # Karafka.monitor.subscribe(Karafka::Instrumentation::ProctitleListener.new)

  routes.draw do
    # Uncomment this if you use Karafka with ActiveJob
    # You ned to define the topic per each queue name you use
    # active_job_topic :default
    topic 'postgres.public.app_demand_profiles' do
      consumer AppDemandProfilesConsumer
    end
    topic 'postgres.public.app_mmp_profiles' do
      consumer AppMmpProfilesConsumer
    end
    topic 'postgres.public.apps' do
      consumer AppsConsumer
    end
    topic 'postgres.public.auction_configurations' do
      consumer AuctionConfigurationsConsumer
    end
    topic 'postgres.public.countries' do
      consumer CountriesConsumer
    end
    topic 'postgres.public.demand_source_accounts' do
      consumer DemandSourceAccountsConsumer
    end
    topic 'postgres.public.demand_sources' do
      consumer DemandSourcesConsumer
    end
    topic 'postgres.public.line_items' do
      consumer LineItemsConsumer
    end
    topic 'postgres.public.users' do
      consumer UsersConsumer
    end
  end
end
