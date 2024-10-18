use crate::client::Client;
use crate::Api;
use crate::{
    models, ApiError, GetAuctionResponse, GetConfigResponse, GetOpenApiSpecResponse,
    PostClickResponse, PostLossResponse, PostRewardResponse, PostShowResponse, PostStatsResponse,
    PostWinResponse,
};
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
    async fn get_config(
        &self,
        config_request: models::ConfigRequest,
        context: &BidonContext,
    ) -> Result<GetConfigResponse, ApiError> {
        self.client.get_config(config_request, context).await
    }

    async fn get_open_api_spec(
        &self,
        context: &BidonContext,
    ) -> Result<GetOpenApiSpecResponse, ApiError> {
        self.client.get_open_api_spec(context).await
    }

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

    async fn post_click(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        click_request: models::ClickRequest,
        context: &BidonContext,
    ) -> Result<PostClickResponse, ApiError> {
        self.client
            .post_click(ad_type, click_request, context)
            .await
    }

    async fn post_loss(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        loss_request: models::LossRequest,
        context: &BidonContext,
    ) -> Result<PostLossResponse, ApiError> {
        self.client.post_loss(ad_type, loss_request, context).await
    }

    async fn post_reward(
        &self,
        ad_type: models::PostRewardAdTypeParameter,
        reward_request: models::RewardRequest,
        context: &BidonContext,
    ) -> Result<PostRewardResponse, ApiError> {
        self.client
            .post_reward(ad_type, reward_request, context)
            .await
    }

    async fn post_show(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        show_request: models::ShowRequest,
        context: &BidonContext,
    ) -> Result<PostShowResponse, ApiError> {
        self.client.post_show(ad_type, show_request, context).await
    }

    async fn post_stats(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        stats_request: models::StatsRequest,
        context: &BidonContext,
    ) -> Result<PostStatsResponse, ApiError> {
        self.client
            .post_stats(ad_type, stats_request, context)
            .await
    }

    async fn post_win(
        &self,
        ad_type: models::GetAuctionAdTypeParameter,
        win_request: models::WinRequest,
        context: &BidonContext,
    ) -> Result<PostWinResponse, ApiError> {
        self.client.post_win(ad_type, win_request, context).await
    }
}
