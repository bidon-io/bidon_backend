# frozen_string_literal: true

# {
#   "id": "843e6bdc-2a81-4126-be19-da1a023178e0",
#   "test": 1,
#   "at": 1,
#   "tmax": 3000,
#   "app": {
#     "ver": "1.0",
#     "bundle":"org.bidon.demo",
#     "cat": ["IAB24"],
#     "id": "1"
#   },
#   "device": {
#     "ifa": "UUID",
#     "ip": "",
#     "carrier": "",
#     "language": "en",
#     "make": "Google",
#     "model": "Google Pixel 7 Pro",
#     "ua": "Mozilla\/5.0 ...",
#     "pxratio": 2.625,
#     "os": "iOS",
#     "devicetype": "4",
#     "osv": "13",
#     "connectiontype": 2,
#     "js": 1,
#     "h": 2179,
#     "w": 1080,
#     "geo": {}
#   },
#   "imp": [{
#     "id": "9987d8b9-2958-4371-99ec-b545bd0d7a9e",
#     "instl": 0,
#     "secure": 1,
#     "exp": 14400,
#     "bidfloor": 30,
#     "banner": {
#       "w": 320,
#       "h": 50
#     }
#   }],
#   "user": {
#     "data": [{
#       "id": "1",
#       "name": "Bidon",
#       "segment": [{
#         "signal": "ChMKCkJpZE1hY2..."
#       }]
#     }]
#   },
#   "regs": {
#     "gdpr": 0,
#     "consent": "0",
#     "ccpa": 0
#   }
# }

# {
#   "id": "843e6bdc-2a81-4126-be19-da1a023178e0",
#   "seatbid": [
#     {
#       "bid": [
#         {
#           "id": "adba7bca-d172-42b2-9f90-1bff8685fdd1",
#           "impid": "9987d8b9-2958-4371-99ec-b545bd0d7a9e",
#           "price": 50.0,
#           "adid": "1378ygfvn928ouyghf19oiuhg03r",
#           "nurl": "https://...",
#           "burl": "https://...",
#           "lurl": "https://...",
#           "adomain": [
#             "bidmachine.io"
#           ],
#           "cid": "phone_banner",
#           "crid": "phone_banner",
#           "h": 50,
#           "w": 320,
#           "ext": {
#             "signaldata": "...AQUIwAIQMg=="
#           }
#         }
#       ],
#       "seat": "3",
#       "group": 0
#     }
#   ],
#   "cur": "USD"
# }

module Api
  module Bidding
    module Demand
      class BidMachine
        ENDPOINT = URI('https://api-eu.bidmachine.io/auction/prebid/applovin')
        HEADERS  = { 'Content-Type' => 'application/json' }.freeze

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

        attr_reader :request, :token, :bidfloor

        delegate :params, to: :request

        def initialize(request, token, bidfloor)
          @request = request
          @token = token
          @bidfloor = bidfloor
        end

        # @return [DemandResponse] Encoded JSON request, Encoded JSON response, response status, price, seatbid
        def call
          data = build_request_body.to_json
          response = Net::HTTP.post(ENDPOINT, data, HEADERS)

          return empty_response(data, response) if response.code_type == Net::HTTPNoContent
          return error_response(data, response) if response.error_type == Net::HTTPClientException

          success_response(data, response)
        end

        private

        def build_request_body # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
          data = {
            id:     SecureRandom.uuid,
            test:   params[:test] ? 1 : 0,
            at:     1,
            tmax:   5000,
            app:    {
              ver:    params[:app][:version],
              bundle: params[:app][:bundle],
              id:     '1',
            },
            user:   {
              data: [user],
            },
            device: params[:device],
            imp:    [imp],
            regs:   {
              coppa: params[:regs][:coppa] ? 1 : 0,
              gdpr:  params[:regs][:gdpr] ? 1 : 0,
            },
          }

          apply_overrides!(data)

          data
        end

        def imp
          res = {
            id:                SecureRandom.uuid,
            displaymanager:    'BidMachine',
            displaymanagerver: params[:adapters][:bidmachine][:sdk_version],
            secure:            1,
            bidfloor:,
          }

          res.merge(ad_type_params)
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

        def user
          {
            id:      '1',
            name:    'Bidon',
            segment: [{ signal: token }],
          }
        end

        def apply_overrides!(data)
          # accuracy -> Int
          if (accuracy = data.dig(:device, :geo, :accuracy))
            data[:device][:geo][:accuracy] = accuracy.round
          end
          # lastfix -> Int, seconds after last geo retrieval, we have unix timestamp of last retrieval
          if (lastfix = data.dig(:device, :geo, :lastfix)) # rubocop:disable Style/GuardClause
            data[:device][:geo][:lastfix] = (Time.zone.now - Time.zone.at(lastfix / 1000)).round
          end
        end

        def parse_bid(bid)
          {
            id:        bid['id'],
            impid:     bid['impid'],
            price:     bid['price'], # Bid price expressed as CPM
            payload:   bid['ext']['signaldata'],
            demand_id: 'bidmachine',
          }
        end

        def empty_response(request, response)
          DemandResponse.new(
            demand:       'bidmachine',
            raw_request:  request,
            raw_response: '',
            status:       response.code,
            price:        0,
            bid:          {},
          )
        end

        def error_response(request, response)
          DemandResponse.new(
            demand:       'bidmachine',
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
            demand:       'bidmachine',
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
