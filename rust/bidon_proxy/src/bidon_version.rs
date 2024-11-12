use hyper;
use std::fmt;

// TODO use the following constants in your code

/// Header - `x-bidon-version` - version of the bidon server.
pub const X_BIDON_VERSION_HEADER: &str = "x-bidon-version";

pub const BIDON_VERSION: &str = "0.0.1"; // TODO: Update version

/// Wrapper for a string being used as an X-Span-ID.
#[derive(Debug, Clone)]
pub struct XBidonVersionString(pub String);

impl XBidonVersionString {
    /// Extract an X-Bidon-version from a request header if present, and if not
    /// generate a new one.
    pub fn get_or_generate<T>(req: &hyper::Request<T>) -> Self {
        let x_bidon_version = req.headers().get(X_BIDON_VERSION_HEADER);

        x_bidon_version
            .and_then(|x| x.to_str().ok())
            .map(|x| XBidonVersionString(x.to_string()))
            .unwrap_or_default()
    }
}

impl Default for XBidonVersionString {
    fn default() -> Self {
        XBidonVersionString(BIDON_VERSION.to_string())
    }
}

impl fmt::Display for XBidonVersionString {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}
