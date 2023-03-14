# frozen_string_literal: true

require 'rails_helper'

RSpec.describe Api::Auction::RoundsFilterer do
  context 'adapter is present only in round 1' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_INTERSTITIAL_1',
          'demands' => %w[admob bidmachine],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_INTERSTITIAL_2',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
      ]
    end

    let(:adapters) do
      {
        'admob' => {
          'version'     => '0.1.0.1-Beta',
          'sdk_version' => '21.5.0',
        },
      }
    end

    let(:expected_result) do
      [
        {
          'id'      => 'ROUND_INTERSTITIAL_1',
          'demands' => %w[admob],
          'timeout' => 15_000,
        },
      ]
    end

    it 'returns filtered rounds in 1 round' do
      expect(described_class.new(rounds:, adapters:).fetch).to eq expected_result
    end
  end

  context 'adapter is not present in rounds' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_INTERSTITIAL_1',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_INTERSTITIAL_2',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
      ]
    end

    let(:adapters) do
      {
        'admob' => {
          'version'     => '0.1.0.1-Beta',
          'sdk_version' => '21.5.0',
        },
      }
    end

    let(:expected_result) { [] }

    it 'returns empty array' do
      expect(described_class.new(rounds:, adapters:).fetch).to eq expected_result
    end
  end

  context 'no adapters passed' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_INTERSTITIAL_1',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_INTERSTITIAL_2',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
      ]
    end

    let(:adapters) do
      {}
    end

    let(:expected_result) { [] }

    it 'returns empty array' do
      expect(described_class.new(rounds:, adapters:).fetch).to eq expected_result
    end
  end

  context 'adapter is present only in rounds 1 and 3' do
    let(:rounds) do
      [
        {
          'id'      => 'ROUND_INTERSTITIAL_1',
          'demands' => %w[admob bidmachine],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_INTERSTITIAL_2',
          'demands' => %w[bidmachine],
          'timeout' => 15_000,
        },
        {
          'id'      => 'ROUND_INTERSTITIAL_3',
          'demands' => %w[admob bidmachine],
          'timeout' => 15_000,
        },
      ]
    end

    let(:adapters) do
      {
        'admob' => {
          'version'     => '0.1.0.1-Beta',
          'sdk_version' => '21.5.0',
        },
      }
    end

    let(:expected_result) do
      [
        {
          'id'      => 'ROUND_INTERSTITIAL_1',
          'demands' => %w[admob],
          'timeout' => 15_000,
        },
      ]
    end

    it 'returns filtered rounds in 1 round' do
      expect(described_class.new(rounds:, adapters:).fetch).to eq expected_result
    end
  end
end
