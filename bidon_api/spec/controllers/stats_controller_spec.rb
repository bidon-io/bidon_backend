# frozen_string_literal: true

require 'rails_helper'

RSpec.describe StatsController, type: :controller do
  context 'X-BidOn-Version header present' do
    before do
      request.headers['X-BidOn-Version'] = '1.2.3'
    end

    context 'valid request' do
      let(:expected_response) do
        { success: true }.to_json
      end

      it 'returns 200 with ok' do
        allow_any_instance_of(Api::Request).to receive(:valid?).and_return(true)
        post :create, params: stats_params.merge(ad_type: 'banner'), as: :json

        expect(response).to have_http_status(:ok)
        expect(response.body).to eq expected_response
      end
    end
  end
end
