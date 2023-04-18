# frozen_string_literal: true

require 'rails_helper'

RSpec.describe Api::Auction::RoundsFilterer do
  subject(:filterer) { described_class.new(rounds:, adapters:) }

  let(:adapters) do
    {
      'unityads' => {
        'version'     => '0.1.1.0-Beta',
        'sdk_version' => '4.5.0',
      },
    }
  end

  context 'when round has networks that are not present in adapters' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_BANNER_2',
          'demands' => %w[admob unityads applovin],
          'timeout' => 15_000,
        },
      ]
    end

    it 'filters not present networks from this round' do
      expect(filterer.fetch).to eq(
        [
          {
            'id'      => 'ROUND_BANNER_2',
            'demands' => %w[unityads],
            'timeout' => 15_000,
          },
        ],
      )
    end
  end

  context 'when round does not have any network from adapters' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_BANNER_3',
          'demands' => %w[dtexchange bidmachine],
          'timeout' => 15_000,
        },
      ]
    end

    it 'does not return this round' do
      expect(filterer.fetch).to eq([])
    end
  end

  context 'when configuration has multiple rounds' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_BANNER_1',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_BANNER_2',
          'demands' => %w[admob unityads applovin],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_BANNER_3',
          'demands' => %w[dtexchange bidmachine unityads],
          'timeout' => 15_000,
        },
      ]
    end

    it 'filters each round independently' do
      expect(filterer.fetch).to eq(
        [
          {
            'id'      => 'ROUND_BANNER_2',
            'demands' => %w[unityads],
            'timeout' => 15_000,
          },
          {
            'id'      => 'ROUND_BANNER_3',
            'demands' => %w[unityads],
            'timeout' => 15_000,
          },

        ],
      )
    end
  end
end
