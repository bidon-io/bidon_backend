require 'rails_helper'

RSpec.describe User, type: :model do
  it 'checks validation' do
    expect(described_class.new).to be_invalid
  end
end
