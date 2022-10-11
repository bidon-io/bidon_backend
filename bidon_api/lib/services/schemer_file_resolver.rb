class SchemerFileResolver
  WINDOWS_URI_PATH_REGEX = %r{\A/[a-z]:}i

  def call(uri)
    path = uri.path
    path = path[1..] if path.match?(WINDOWS_URI_PATH_REGEX)
    JSON.parse(File.read(URI::DEFAULT_PARSER.unescape(path)))
  end
end
