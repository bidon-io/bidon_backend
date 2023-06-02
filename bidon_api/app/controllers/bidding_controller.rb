# frozen_string_literal: true

class BiddingController < ApplicationController
  def create
    auction_response = Api::Bidding::Response.new(api_request)

    if auction_response.bid?
      render json: auction_response.body, status: :ok
    else
      render nothing: true, status: :no_content
    end
  end

  private

  def schema_file_name
    'bidding.json'
  end
end
