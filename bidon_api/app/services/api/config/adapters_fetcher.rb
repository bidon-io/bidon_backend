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

      private

      def fetch_adapter(adapter_name) # rubocop:disable Metrics/MethodLength
        return {} unless app_mmp_profile

        case adapter_name.to_s
        when 'appsflyer'
          {
            dev_key: app_mmp_profile.appsflyer_dev_key,
            app_id:  app_mmp_profile.appsflyer_app_id,
          }
        when 'adjust'
          {
            app_token: app_mmp_profile.adjust_app_token,
            s2s_token: app_mmp_profile.adjust_s2s_token,
          }
        when 'bidmachine'
          {
            seller_id:        '1',
            endpoint:         'x.appbaqend.com',
            mediation_config: %w[meta_audience criteo pangle amazon adcolony my_target vungle tapjoy notsy],
          }
        when 'applovin'
          {
            app_key: 'example',
          }
        else
          {}
        end
      end
    end
  end
end
