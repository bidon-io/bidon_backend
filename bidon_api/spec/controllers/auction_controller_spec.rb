# frozen_string_literal: true

require 'rails_helper'

RSpec.describe AuctionController, type: :controller do
  context 'missing X-BidOn-Version header' do
    let(:expected_response) do
      {
        error: {
          code:    422,
          message: 'Request should contain X-BidOn-Version header',
        },
      }.to_json
    end

    it 'returns 422 with error' do
      post :create, params: auction_params.merge(ad_type: 'banner'), as: :json

      expect(response).to have_http_status(:unprocessable_entity)
      expect(response.body).to eq expected_response
    end
  end

  context 'X-BidOn-Version header present' do
    before do
      request.headers['X-BidOn-Version'] = '1.2.3'
    end

    context 'valid response' do
      let(:expected_response) do
        {
          'rounds'     => [],
          'line_items' => [],
          'token'      => '{}',
          'pricefloor' => 0,
        }.to_json
      end

      it 'returns 200 with ok' do
        allow_any_instance_of(Api::Request).to receive(:valid?).and_return(true)
        allow_any_instance_of(Api::Auction::Response).to receive(:body).and_return(expected_response)

        post :create, params: auction_params.merge(ad_type: 'banner'), as: :json

        expect(response).to have_http_status(:ok)
        expect(response.body).to eq expected_response
      end
    end

    context 'empty response' do
      let(:expected_response) do
        {
          error: {
            code:    422,
            message: 'No ads found',
          },
        }.to_json
      end

      it 'returns 422 with No ads error' do
        allow_any_instance_of(Api::Request).to receive(:valid?).and_return(true)
        allow_any_instance_of(Api::Auction::Response).to receive(:present?).and_return(false)

        post :create, params: auction_params.merge(ad_type: 'banner'), as: :json

        expect(response).to have_http_status(:unprocessable_entity)
        expect(response.body).to eq expected_response
      end
    end

    context 'invalid request' do
      let(:expected_response) do
        {
          error: {
            code:    422,
            message: 'App key is invalid',
          },
        }.to_json
      end

      it 'returns 422 with error' do
        allow_any_instance_of(Api::Request).to receive(:valid?).and_return(false)

        post :create, params: auction_params.merge(ad_type: 'banner'), as: :json

        expect(response).to have_http_status(:unprocessable_entity)
        expect(response.body).to eq expected_response
      end
    end

    context 'error request' do
      let(:expected_response) do
        {
          error: {
            code:    500,
            message: 'Internal Server Error',
          },
        }.to_json
      end

      it 'returns 500 with error' do
        allow(Api::Request).to receive(:new).and_raise(StandardError)

        post :create, params: auction_params.merge(ad_type: 'banner'), as: :json

        expect(response).to have_http_status(:internal_server_error)
        expect(response.body).to eq expected_response
      end
    end
  end
end
