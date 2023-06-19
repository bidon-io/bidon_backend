# frozen_string_literal: true

module Api
  module Bidding
    module Demand
      class BidMachine
        CONFIG = {
          endpoint:  URI(Utils.fetch_from_env('BIDDING_BIDMACHINE_URL')),
          demand_id: 'bidmachine',
        }.freeze

        HEADERS = { 'Content-Type' => 'application/json' }.freeze

        BANNER_FORMATS = {
          'BANNER'      => [320, 50],
          'LEADERBOARD' => [728, 90],
          'MREC'        => [300, 250],
          'ADAPTIVE'    => [0, 50],
        }.freeze

        FULLSCREEN_FORMATS = {
          'PHONE'  => [320, 480],
          'TABLET' => [768, 1024],
        }.freeze

        attr_reader :request, :ip, :token, :bidfloor

        delegate :params, :app, to: :request

        def initialize(request, ip, token, bidfloor)
          @request = request
          @ip = ip
          @token = token
          @bidfloor = bidfloor
        end

        # @return [DemandResponse] Encoded JSON request, Encoded JSON response, response status, price, seatbid
        def call
          data = build_request_body.to_json
          response = Net::HTTP.post(CONFIG[:endpoint], data, HEADERS)

          return empty_response(data, response) if response.code_type  == Net::HTTPNoContent
          return error_response(data, response) if response.error_type == Net::HTTPClientException

          success_response(data, response)
        end

        private

        def build_request_body # rubocop:disable Metrics/AbcSize
          {
            id:     SecureRandom.uuid,
            test:   params[:test] ? 1 : 0,
            at:     1,
            tmax:   5000,
            app:    {
              ver:       params.dig(:app, :version).to_s,
              bundle:    app.package_name,
              id:        app.id.to_s,
              publisher: { id: adapter_config[:seller_id].to_s },
            },
            device:,
            imp:    [imp],
            regs:   {
              coppa: params.dig(:regs, :coppa) ? 1 : 0,
              gdpr:  params.dig(:regs, :gdpr)  ? 1 : 0,
            },
          }
        end

        def imp
          res = {
            id:                SecureRandom.uuid,
            displaymanager:    CONFIG[:demand_id],
            displaymanagerver: params.dig(:adapters, :bidmachine, :sdk_version),
            secure:            1,
            bidfloor:,
            ext:               { bid_token: token },
          }

          res.deep_merge(ad_type_params)
        end

        def ad_type_params # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
          if params[:imp].key?(:banner)
            size = BANNER_FORMATS[params.dig(:imp, :banner, :format)]
            size.reverse! if params.dig(:imp, :orientation) == 'LANDSCAPE'

            { instl:  0,
              banner: {
                w:     size[0],
                h:     size[1],
                btype: [],
                battr: [1, 2, 5, 8, 9, 14, 17],
                pos:   1,
              } }
          elsif params[:imp].key?(:interstitial)
            size = FULLSCREEN_FORMATS[params.dig(:device, :type)]
            size.reverse! if params.dig(:imp, :orientation) == 'LANDSCAPE'

            { instl:  1,
              banner: {
                w:     size[0],
                h:     size[1],
                btype: [],
                battr: [],
                pos:   7,
              } }
          elsif params[:imp].key?(:rewarded)
            size = FULLSCREEN_FORMATS[params.dig(:device, :type)]
            size.reverse! if params.dig(:imp, :orientation) == 'LANDSCAPE'

            { instl:  1,
              ext:    { rewarded: 1 },
              banner: {
                w:     320,
                h:     480,
                btype: [],
                battr: [16],
                pos:   7,
              } }
          else
            {}
          end
        end

        def device
          Bidding::DeviceBuilder.new(params[:device], params[:user], ip).call
        end

        def adapter_config
          @adapter_config ||= Api::Config::AdaptersFetcher.new(
            app:, config_adapters: request.adapters,
          ).fetch_adapter(CONFIG[:demand_id])
        end

        def parse_bid(bid)
          {
            id:        bid['id'],
            impid:     bid['impid'],
            price:     bid['price'], # Bid price expressed as CPM
            payload:   bid['adm'],
            demand_id: CONFIG[:demand_id],
          }
        end

        def empty_response(request, response)
          DemandResponse.new(
            demand:       CONFIG[:demand_id],
            raw_request:  request,
            raw_response: '',
            status:       response.code,
            price:        0,
            bid:          {},
          )
        end

        def error_response(request, response)
          DemandResponse.new(
            demand:       CONFIG[:demand_id],
            raw_request:  request,
            raw_response: { error: response.body }.to_json,
            status:       response.code,
            price:        0,
            bid:          {},
          )
        end

        def success_response(request, response)
          bid = JSON.parse(response.body)['seatbid'][0]['bid'][0]

          DemandResponse.new(
            demand:       CONFIG[:demand_id],
            raw_request:  request,
            raw_response: response.body,
            status:       response.code,
            price:        bid['price'],
            bid:          parse_bid(bid),
          )
        end
      end
    end
  end
end
