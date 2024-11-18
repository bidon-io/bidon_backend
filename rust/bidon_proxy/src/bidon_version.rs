use axum::body::BoxBody;
use axum::{
    body::{boxed, Body},
    extract::Extension,
    http::{HeaderValue, Request},
    middleware::{self, Next},
    response::Response,
    Router,
};
use hyper::StatusCode;
use std::fmt;
use std::net::SocketAddr;
// TODO use the following constants in your code

/// Header - `x-bidon-version` - version of the bidon server.
pub const X_BIDON_VERSION_HEADER: &str = "x-bidon-version";

pub const BIDON_VERSION: &str = "0.0.1"; // TODO: Update version

/// Wrapper for a string being used as an X-Span-ID.
#[derive(Debug, Clone)]
pub struct XBidonVersionString(pub String);

impl XBidonVersionString {
    pub async fn extract_header_middleware<B>(
        mut req: Request<B>,
        next: Next<B>,
    ) -> Response<BoxBody>
    where
        B: Send + 'static,
    {
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
            let response = next.run(req).await;
            response.map(boxed)
        } else {
            // Return an error response if the header is missing
            let body = boxed(Body::from(format!(
                "Missing '{}' header",
                X_BIDON_VERSION_HEADER
            )));
            Response::builder()
                .status(StatusCode::BAD_REQUEST)
                .body(body)
                .unwrap()
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
