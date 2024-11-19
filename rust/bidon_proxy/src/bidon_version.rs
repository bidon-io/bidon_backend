use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::{
    body::Body,
    extract::{Extension, Request},
    http::HeaderValue,
    middleware::{self, Next},
    response::Response,
    // response::IntoResponse,
    Router,
};
use std::fmt;
use std::net::SocketAddr;
use tonic::IntoRequest;
// TODO use the following constants in your code

/// Header - `x-bidon-version` - version of the bidon server.
pub const X_BIDON_VERSION_HEADER: &str = "x-bidon-version";

pub const BIDON_VERSION: &str = "0.0.1"; // TODO: Update version

/// Wrapper for a string being used as an X-Span-ID.
#[derive(Debug, Clone)]
pub struct XBidonVersionString(pub String);

impl XBidonVersionString {
    // #[axum::debug_middleware]
    pub async fn extract_header_middleware(
        mut req: axum::extract::Request,
        next: Next,
    ) -> Response {
        if let Some(bidon_version) = req
            .headers()
            .get(X_BIDON_VERSION_HEADER)
            .and_then(|x| x.to_str().ok())
            .map(|x| XBidonVersionString(x.to_string()))
        {
            // Store the header value in the request extensions
            req.extensions_mut().insert(bidon_version);
            // Continue processing the request
            // Run the next middleware or handler, and ensure the response body is boxed
            next.run(req).await
        } else {
            // Return an error response if the header is missing
            (
                StatusCode::BAD_REQUEST,
                format!("Missing {} header", X_BIDON_VERSION_HEADER),
            )
                .into_response()
        }
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
