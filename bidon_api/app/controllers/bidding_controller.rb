# frozen_string_literal: true

class BiddingController < ApplicationController
  def create
    auction_response = Api::Bidding::Response.new(api_request)

    render json: auction_response.body, status: :ok
  end

  private

  def schema_file_name
    'bidding.json'
  end
end
