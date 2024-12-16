use axum::async_trait;
use axum::extract::rejection::JsonRejection;
use axum::extract::{FromRequest, FromRequestParts, Json, Path, Request};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum_client_ip::InsecureClientIp;

use crate::adapter;
use crate::com::iabtechlab::openrtb::v3::Openrtb;
use crate::sdk::{AuctionRequest, GetAuctionAdTypeParameter};

pub struct AuctionRequestPayload(pub Openrtb);

pub enum AuctionRequestRejection {
    MissingBidonVersionHeader,
    MissingIpAddress,
    InvalidJson(JsonRejection),
    InvalidBiddingRequest,
    InvalidAdType,
}

impl IntoResponse for AuctionRequestRejection {
    fn into_response(self) -> Response {
        match self {
            AuctionRequestRejection::InvalidJson(json_error) => (
                StatusCode::BAD_REQUEST,
                format!("Invalid JSON in request body: {:?}", json_error),
            )
                .into_response(),
            AuctionRequestRejection::InvalidBiddingRequest => {
                (StatusCode::BAD_REQUEST, "Invalid bidding request").into_response()
            }
            AuctionRequestRejection::MissingBidonVersionHeader => {
                (StatusCode::BAD_REQUEST, "Missing x-bidon-version header").into_response()
            }
            AuctionRequestRejection::MissingIpAddress => {
                (StatusCode::BAD_REQUEST, "Missing client ip address").into_response()
            }
            AuctionRequestRejection::InvalidAdType => {
                (StatusCode::BAD_REQUEST, "Invalid ad type").into_response()
            }
        }
    }
}

#[async_trait]
impl<B> FromRequest<B> for AuctionRequestPayload
where
    B: Send + Sync + 'static,
{
    type Rejection = AuctionRequestRejection;

    async fn from_request(req: Request, state: &B) -> Result<Self, Self::Rejection> {
        /// Header - `x-bidon-version` - version of the bidon server.
        const X_BIDON_VERSION_HEADER: &str = "x-bidon-version";

        let bidon_version = req
            .headers()
            .get(X_BIDON_VERSION_HEADER)
            .and_then(|x| x.to_str().ok())
            .map(|bidon_version| bidon_version.to_string())
            // Return an error response if the header is missing
            .ok_or(AuctionRequestRejection::InvalidBiddingRequest)?;

        let (mut parts, body) = req.into_parts();
        let InsecureClientIp(ip) = InsecureClientIp::from(&parts.headers, &parts.extensions)
            .map_err(|_| AuctionRequestRejection::MissingIpAddress)?;

        let Path(ad_type) =
            Path::<GetAuctionAdTypeParameter>::from_request_parts(&mut parts, state)
                .await
                .map_err(|_| AuctionRequestRejection::InvalidAdType)?;

        let Json(auction_request) =
            Json::<AuctionRequest>::from_request(Request::from_parts(parts, body), state)
                .await
                .map_err(AuctionRequestRejection::InvalidJson)?;

        let openrtb_request = adapter::try_from(&auction_request, bidon_version, ip, ad_type)
            .map_err(|_| AuctionRequestRejection::InvalidBiddingRequest)?;

        Ok(AuctionRequestPayload(openrtb_request))
    }
}
