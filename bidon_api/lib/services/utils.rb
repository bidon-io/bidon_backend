module Utils
  module_function

  # @param [Hash] params
  #
  # @return [String]
  def encode_params(params)
    Base64.encode64(zip(params.to_json))
  end

  # @param [String] raw_params
  #
  # @return [String]
  def decode_params(raw_params)
    unzip(Base64.decode64(raw_params))
  end

  # @param [String] data
  #
  # @return [String]
  def zip(data)
    ActiveSupport::Gzip.compress(data)
  end

  # @param [String] data
  #
  # @return [String]
  def unzip(data)
    ActiveSupport::Gzip.decompress(data)
  end

  # @param [Hash] hash
  #
  # @return [Hash]
  def smash_hash(hash)
    Flatten.smash(hash, smash_array: true, separator: '__')
  end

  # @param [String] key
  #
  # @raise [KeyError]
  # @return [String]
  def fetch_from_env(key)
    ENV.fetch(key) { raise KeyError, "Missing #{key} env variable" }
  end
end
