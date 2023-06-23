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

      def app_demand_profile(adapter)
        AppDemandProfile
          .eager(:demand_source_account)
          .where(app_id: app.id, account_type: AppDemandProfile::KEY_TO_ACCOUNT_TYPE[adapter])
          .first
      end

      def fetch_adapter(adapter_name)
        case adapter_name.to_s
        when 'appsflyer'
          fetch_appsflyer_adapter
        when 'adjust'
          fetch_adjust_adapter
        when *AppDemandProfile::ADAPTERS_LIST
          fetch_demand_adapter(adapter_name.to_s)
        else
          {}
        end
      end

      private

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

      def fetch_demand_adapter(adapter)
        profile = app_demand_profile(adapter)
        return {} unless profile

        extra = JSON.parse(profile.demand_source_account.extra)

        case adapter
        when 'applovin'
          { app_key: extra['api_key'] }
        when 'bidmachine'
          {
            seller_id:        extra['seller_id'],
            endpoint:         extra['endpoint'],
            mediation_config: extra['mediation_config'],
          }
        else
          extra
        end
      end
    end
  end
end
