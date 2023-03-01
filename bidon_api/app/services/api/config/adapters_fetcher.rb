module Api
  module Config
    class AdaptersFetcher
      prepend MemoWise

      attr_reader :app, :config_adapters

      def initialize(app:, config_adapters:)
        @app = app
        @config_adapters = config_adapters
      end

      def fetch
        config_adapters.each_with_object({}) do |(adapter_name, _v), result|
          result[adapter_name] = fetch_adapter(adapter_name)
        end
      end

      def app_mmp_profile
        AppMmpProfile.where(app_id: app.id).order(Sequel.desc(:start_date)).first
      end
      memo_wise :app_mmp_profile

      def applovin_demand_profile
        AppDemandProfile.eager(:demand_source_account)
                        .where(app_id: app.id, account_type: 'DemandSourceAccount::Applovin').first
      end
      memo_wise :applovin_demand_profile

      def bidmachine_demand_profile
        AppDemandProfile.eager(:demand_source_account)
                        .where(app_id: app.id, account_type: 'DemandSourceAccount::BidMachine').first
      end
      memo_wise :bidmachine_demand_profile

      def data_exchange_demand_profile
        AppDemandProfile.eager(:demand_source_account)
                        .where(app_id: app.id, account_type: 'DemandSourceAccount::DataExchange').first
      end
      memo_wise :data_exchange_demand_profile

      def unity_ads_demand_profile
        AppDemandProfile.eager(:demand_source_account)
                        .where(app_id: app.id, account_type: 'DemandSourceAccount::UnityAds').first
      end
      memo_wise :unity_ads_demand_profile

      private

      def fetch_adapter(adapter_name)
        case adapter_name.to_s
        when 'appsflyer'
          fetch_appsflyer_adapter
        when 'adjust'
          fetch_adjust_adapter
        when 'bidmachine'
          fetch_bidmachine_adapter
        when 'applovin'
          fetch_applovin_adapter
        when 'dtexchange'
          fetch_data_exchange_adapter
        when 'unityads'
          fetch_unity_ads_adapter
        else
          {}
        end
      end

      def fetch_appsflyer_adapter
        return {} unless app_mmp_profile

        {
          dev_key: app_mmp_profile.appsflyer_dev_key,
          app_id:  app_mmp_profile.appsflyer_app_id,
        }
      end

      def fetch_adjust_adapter
        return {} unless app_mmp_profile

        {
          app_token: app_mmp_profile.adjust_app_token,
          s2s_token: app_mmp_profile.adjust_s2s_token,
        }
      end

      def fetch_applovin_adapter
        return {} unless applovin_demand_profile

        extra = JSON.parse(applovin_demand_profile.demand_source_account.extra)

        {
          app_key: extra['api_key'],
        }
      end

      def fetch_bidmachine_adapter
        return {} unless bidmachine_demand_profile

        extra = JSON.parse(bidmachine_demand_profile.demand_source_account.extra)

        {
          seller_id:        extra['seller_id'],
          endpoint:         extra['endpoint'],
          mediation_config: extra['mediation_config'],
        }
      end

      def fetch_data_exchange_adapter
        return {} unless data_exchange_demand_profile

        JSON.parse(data_exchange_demand_profile.demand_source_account.extra)
      end

      def fetch_unity_ads_adapter
        return {} unless unity_ads_demand_profile

        JSON.parse(unity_ads_demand_profile.demand_source_account.extra)
      end
    end
  end
end
