# frozen_string_literal: true

module Api
  module Bidding
    class DeviceBuilder
      attr_reader :device_params, :user_params, :ip

      # https://github.com/InteractiveAdvertisingBureau/AdCOM/blob/master/AdCOM%20v1.0%20FINAL.md#list--connection-types-
      CONNECTION_TYPE_MAP = {
        'ETHERNET'         => 1,
        'WIFI'             => 2,
        'CELLULAR'         => 3,
        'CELLULAR_UNKNOWN' => 3,
        'CELLULAR_2_G'     => 4,
        'CELLULAR_3_G'     => 5,
        'CELLULAR_4_G'     => 6,
        'CELLULAR_5_G'     => 7,
      }.freeze
      DEFAULT_CONNECTION_TYPE = 3

      # https://github.com/InteractiveAdvertisingBureau/AdCOM/blob/master/AdCOM%20v1.0%20FINAL.md#list--device-types-
      DEVICE_TYPE_MAP = {
        'PHONE' => 4, 'TABLET' => 5
      }.freeze
      DEFAULT_DEVICE_TYPE = 4

      def initialize(device_params, user_params, ip)
        @device_params = device_params
        @user_params = user_params
        @ip = ip
      end

      def call # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
        data = {
          ip:,
          w:              device_params[:w],
          h:              device_params[:h],
          js:             device_params[:js],
          devicetype:     DEVICE_TYPE_MAP.fetch(device_params[:type], DEFAULT_DEVICE_TYPE),
          connectiontype: CONNECTION_TYPE_MAP.fetch(device_params[:connection_type], DEFAULT_CONNECTION_TYPE),
          os:             device_params[:os],
          osv:            device_params[:osv],
          pxratio:        device_params[:pxratio],
          language:       device_params[:language],
          make:           device_params[:make],
          hwv:            device_params[:hwv],
          ua:             device_params[:ua],
          ppi:            device_params[:ppi],
          model:          device_params[:model],
          ifa:            user_params[:idfa],
        }

        data[:carrier] = device_params[:carrier] if device_params[:carrier].present?
        data[:mccmnc]  = device_params[:mccmnc]  if device_params[:mccmnc].present?
        data[:geo]     = geo_hash unless geocoder.unknown_country

        data.compact
      end

      private

      def geo_hash
        {
          lat:       geocoder.lat,
          lon:       geocoder.lon,
          type:      2,
          ipservice: geocoder.ip_service,
          country:   geocoder.country_code3,
          region:    geocoder.region_code,
          city:      geocoder.city_name,
          zip:       geocoder.zip_code,
          accuracy:  geocoder.accuracy,
        }.compact
      end

      def geocoder
        @geocoder ||= OfflineGeocoder.find_geo_data(ip)
      end
    end
  end
end
