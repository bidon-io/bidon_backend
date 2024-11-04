use crate::client::Client;
use crate::Api;
use crate::{models, ApiError, GetAuctionResponse};
use async_trait::async_trait;
use hyper::header::{HeaderName, HeaderValue, CONTENT_TYPE};
use hyper::{service::Service, Body, Request, Response, Uri};
use std::fmt;

use crate::bidon_version::XBidonVersionString;
use crate::context::BidonContext;
use swagger::{make_context_ty, new_context_type, AuthData, ContextBuilder, EmptyContext, Has, Push, XSpanIdString};

pub struct ProxyServer<S>
where
    S: Service<(Request<Body>, BidonContext), Response=Response<Body>>
    + Clone
    + Sync
    + Send
    + 'static,
    S::Future: Send + 'static,
    S::Error: Into<crate::ServiceError> + fmt::Display,
{
    client: Client<S, BidonContext>,
}

#[allow(dead_code)]
impl<S> ProxyServer<S>
where
    S: Service<(Request<Body>, BidonContext), Response=Response<Body>>
    + Clone
    + Sync
    + Send
    + 'static,
    S::Future: Send + 'static,
    S::Error: Into<crate::ServiceError> + fmt::Display,
{
    pub fn new(client: Client<S, BidonContext>) -> Self {
        ProxyServer { client }
    }
}

#[async_trait]
impl<S> Api<BidonContext> for ProxyServer<S>
where
    S: Service<(Request<Body>, BidonContext), Response=Response<Body>>
    + Clone
    + Sync
    + Send
    + 'static,
    S::Future: Send + 'static,
    S::Error: Into<crate::ServiceError> + fmt::Display,
{
    async fn get_auction(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        auction_request: models::AuctionRequest,
        context: &BidonContext,
    ) -> Result<GetAuctionResponse, ApiError> {
        self.client
            .get_auction(ad_type, auction_request, context)
            .await
    }
}
