# frozen_string_literal: true

class AuctionController < ApplicationController
  def create
    api_request = Api::Request.new(params)

    if api_request.valid?
      auction_response = Api::Auction::Response.new(api_request)

      if auction_response.present?
        render json: auction_response.body, status: :ok
      else
        render_empty_result
      end
    else
      render_app_key_invalid
    end
  end
end
